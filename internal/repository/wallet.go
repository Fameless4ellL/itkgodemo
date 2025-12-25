package repository

import (
	"itkdemo/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(db *gorm.DB) *Postgres {
	return &Postgres{db: db}
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

func (p *Postgres) Update(wallet *domain.Wallet) error {
	return p.db.Model(&wallet).Update("balance", wallet.Balance).Error

}

func (p *Postgres) Delete(id uuid.UUID) error {
	return p.db.Delete(&domain.Wallet{}, "id = ?", id).Error
}
