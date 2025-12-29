package repository

import (
	"context"
	"itkdemo/internal/domain"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var sfGroup singleflight.Group

type Postgres struct {
	db    *gorm.DB
	cache *ristretto.Cache[string, int64]
	task  chan domain.Task
}

func NewPostgres(db *gorm.DB, cache *ristretto.Cache[string, int64]) *Postgres {
	repo := &Postgres{
		db:    db,
		cache: cache,
		task:  make(chan domain.Task, 5000),
	}

	go repo.Run()

	return repo
}

func (p *Postgres) Create(wallet *domain.Wallet) error {
	return p.db.Create(wallet).Error
}

func (p *Postgres) GetByID(id uuid.UUID) (*domain.Wallet, error) {
	if val, found := p.cache.Get(id.String()); found {
		return &domain.Wallet{
			ID:      id,
			Balance: val,
		}, nil
	}

	// avoid spam
	val, err, _ := sfGroup.Do(id.String(), func() (any, error) {
		var wallet domain.Wallet
		if err := p.db.First(&wallet, "id = ?", id).Error; err != nil {
			return 0, err
		}
		p.cache.SetWithTTL(id.String(), wallet.Balance, 1, 5*time.Minute)
		return wallet.Balance, nil
	})

	if err != nil {
		return nil, err
	}
	return &domain.Wallet{
		ID:      id,
		Balance: val.(int64),
	}, nil
}

func (p *Postgres) Update(id uuid.UUID, amount int64) error {
	t := make(chan error, 1)
	task := domain.Task{
		ID:     id,
		Amount: amount,
		Resp:   t,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case p.task <- task:
		return <-t
	case <-ctx.Done():
		return domain.ErrTooBusy
	}
}

func (p *Postgres) Delete(id uuid.UUID) error {
	return p.db.Delete(&domain.Wallet{}, "id = ?", id).Error
}

func (p *Postgres) Run() {
	ticker := time.NewTicker(100 * time.Millisecond)
	batch := make([]domain.Task, 0, 5000)

	for {
		select {
		case task := <-p.task:
			batch = append(batch, task)
			if len(batch) >= 200 {
				p.Batch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				p.Batch(batch)
				batch = batch[:0]
			}
		}
	}
}

func (p *Postgres) Batch(batch []domain.Task) {
	totals := make(map[uuid.UUID]int64)
	for _, task := range batch {
		totals[task.ID] += task.Amount
	}

	err := make(map[uuid.UUID]error)
	for id, amount := range totals {
		res := p.db.
			Model(&domain.Wallet{}).
			Where("id = ? AND (balance + ?) >= 0", id, amount).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount))

		if res.Error != nil {
			err[id] = res.Error
		} else if res.RowsAffected == 0 {
			err[id] = domain.ErrInsufficientBalance
		} else {
			if val, found := p.cache.Get(id.String()); found {
				p.cache.SetWithTTL(id.String(), val+amount, 1, 5*time.Minute)
			}
		}

	}

	for _, task := range batch {
		task.Resp <- err[task.ID]
	}
}
