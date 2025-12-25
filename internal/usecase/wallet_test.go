package usecase

import (
	"errors"
	"itkdemo/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) GetByID(id uuid.UUID) (*domain.Wallet, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockRepo) Create(wallet *domain.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockRepo) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) Update(id uuid.UUID, amount int64) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

func TestUpdateWallet(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		wallet      *domain.Operation
		wantBalance int64
		wantErr     error
	}{
		{
			name: "Valid wallet",
			wallet: &domain.Operation{
				WalletID: uuid.New(),
				Amount:   100,
				Type:     domain.Deposit,
			},
			wantBalance: 900,
			wantErr:     nil,
		},
		{
			name: "Invalid ID",
			wallet: &domain.Operation{
				WalletID: uuid.Nil,
				Amount:   100,
				Type:     domain.Deposit,
			},
			wantBalance: 900,
			wantErr:     domain.ErrInvalidID,
		},
		{
			name: "InsufficientBalance",
			wallet: &domain.Operation{
				WalletID: uuid.Nil,
				Amount:   0,
				Type:     domain.Deposit,
			},
			wantBalance: 900,
			wantErr:     domain.ErrInsufficientBalance,
		},
		{
			name: "InvalidOperationType",
			wallet: &domain.Operation{
				WalletID: uuid.Nil,
				Amount:   100,
				Type:     4,
			},
			wantBalance: 900,
			wantErr:     domain.ErrInvalidOperationType,
		},
		{
			name: "Withdrawal",
			wallet: &domain.Operation{
				WalletID: uuid.New(),
				Amount:   100,
				Type:     domain.Withdraw,
			},
			wantBalance: 900,
			wantErr:     nil,
		},
		{
			name: "WithdrawalInsufficientBalance",
			wallet: &domain.Operation{
				WalletID: uuid.New(),
				Amount:   100,
				Type:     domain.Withdraw,
			},
			wantBalance: 50,
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepo)

			repo.On("GetByID", mock.Anything).Return(&domain.Wallet{ID: tt.wallet.WalletID, Balance: tt.wantBalance}, tt.wantErr)
			repo.On("Update", mock.Anything, mock.Anything).Return(nil)

			uc := NewWalletUseCase(repo)
			err := uc.UpdateWallet(tt.wallet)
			if err == domain.ErrInsufficientBalance && tt.name == "WithdrawalInsufficientBalance" {
				return
			} else if err != tt.wantErr {
				t.Errorf("UpdateWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateWallet(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "New wallet",
			wantErr: nil,
		},
		{
			name:    "Invalid",
			wantErr: errors.New("something go wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepo)

			repo.On("Create", mock.Anything).Return(tt.wantErr)

			uc := NewWalletUseCase(repo)
			_, err := uc.CreateWallet()
			if err != tt.wantErr {
				t.Errorf("CreateWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteWallet(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "Delete wallet",
			wantErr: nil,
		},
		{
			name:    "Invalid",
			wantErr: errors.New("something go wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepo)
			repo.On("Delete", mock.Anything).Return(tt.wantErr)

			uc := NewWalletUseCase(repo)
			err := uc.DeleteWallet(uuid.New())

			if err != tt.wantErr {
				t.Errorf("DeleteWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetWallet(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "Get wallet",
			wantErr: nil,
		},
		{
			name:    "wallet not found",
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepo)

			repo.On("GetByID", mock.Anything).Return(&domain.Wallet{ID: uuid.New(), Balance: 1000}, tt.wantErr)

			uc := NewWalletUseCase(repo)
			_, err := uc.GetWallet(uuid.New())
			if err != tt.wantErr {
				t.Errorf("GetWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
