package usecase

import (
	"itkdemo/internal/domain"

	"github.com/google/uuid"
)

type WalletUseCase struct {
	repo domain.WalletRepository
}

func NewWalletUseCase(repo domain.WalletRepository) *WalletUseCase {
	return &WalletUseCase{
		repo: repo,
	}
}

func (u *WalletUseCase) CreateWallet() (*domain.Wallet, error) {
	wallet := domain.NewWallet()
	if err := u.repo.Create(wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (u *WalletUseCase) DeleteWallet(id uuid.UUID) error {
	if err := u.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (u *WalletUseCase) GetWallet(id uuid.UUID) (*domain.Wallet, error) {
	wallet, err := u.repo.GetByID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return wallet, nil
}

func (u *WalletUseCase) UpdateWallet(op *domain.Operation) error {
	if op.Amount <= 0 {
		return domain.ErrInsufficientBalance
	}
	if op.Type != domain.Deposit && op.Type != domain.Withdraw {
		return domain.ErrInvalidOperationType
	}
	wallet, err := u.repo.GetByID(op.WalletID)
	if err != nil {
		return err
	}
	if op.Type == domain.Withdraw && wallet.Balance < op.Amount {
		return domain.ErrInsufficientBalance
	}
	if op.Type == domain.Withdraw {
		op.Amount = -op.Amount
	}

	return u.repo.Update(wallet.ID, op.Amount)
}
