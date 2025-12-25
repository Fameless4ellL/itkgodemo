package domain

import (
	"github.com/google/uuid"
)

type OperationType int

const (
	Deposit OperationType = iota + 1
	Withdraw
)

type Operation struct {
	WalletID uuid.UUID     `json:"id"`
	Amount   int64         `json:"amount"`
	Type     OperationType `json:"type"`
}

type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance int64     `json:"balance"`
}

func NewWallet() *Wallet {
	return &Wallet{
		ID:      uuid.New(),
		Balance: 1000,
	}
}

type WalletRepository interface {
	Create(wallet *Wallet) error
	GetByID(uuid.UUID) (*Wallet, error)
	Update(wallet *Wallet) error
	Delete(id uuid.UUID) error
}

type WalletUseCase interface {
	CreateWallet() (*Wallet, error)
	GetWallet(id uuid.UUID) (*Wallet, error)
	UpdateWallet(*Operation) error
	DeleteWallet(id uuid.UUID) error
}
