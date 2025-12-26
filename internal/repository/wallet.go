package repository

import (
	"context"
	"itkdemo/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Postgres struct {
	db   *gorm.DB
	task chan domain.Task
}

func NewPostgres(db *gorm.DB) *Postgres {
	repo := &Postgres{
		db:   db,
		task: make(chan domain.Task, 1),
	}

	go repo.Run()

	return repo
}

func (p *Postgres) Create(wallet *domain.Wallet) error {
	return p.db.Create(wallet).Error
}

func (p *Postgres) GetByID(id uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	if err := p.db.First(&wallet, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
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

	err := p.db.Transaction(func(tx *gorm.DB) error {
		for id, amount := range totals {
			if err := tx.Model(&domain.Wallet{ID: id}).Where("balance >= ?", 0).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		for _, task := range batch {
			task.Resp <- err
		}
		return
	}

	for _, task := range batch {
		task.Resp <- nil
	}
}
