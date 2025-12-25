package repository

import (
	"fmt"
	"itkdemo/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) (*Postgres, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	return NewPostgres(gormDB), mock
}

func TestPostgres_CRUD(t *testing.T) {
	repo, mock := SetupTestDB(t)
	walletID := uuid.New()
	wallet := &domain.Wallet{
		ID:      walletID,
		Balance: 100,
	}

	t.Run("CreateWallet", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "wallets"`).
			WithArgs(wallet.ID, wallet.Balance).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(wallet)
		assert.NoError(t, err)
	})

	t.Run("GetWallet", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "balance"}).
			AddRow(walletID, 100)
		mock.ExpectQuery(`SELECT \* FROM "wallets" WHERE .*id = \$1`).
			WithArgs(walletID, sqlmock.AnyArg()).
			WillReturnRows(rows)

		found, err := repo.GetByID(walletID)

		assert.NoError(t, err)
		if assert.NotNil(t, found) {
			assert.Equal(t, walletID, found.ID)
		}
	})

	t.Run("GetWalletErr", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "wallets"`).
			WithArgs(walletID, sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("connection refused"))

		found, err := repo.GetByID(walletID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection refused")
		assert.Nil(t, found)
	})

	t.Run("UpdateWallet", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "wallets" SET`).
			WithArgs(250, walletID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		wallet.Balance = 250
		err := repo.Update(wallet.ID, 1)
		assert.NoError(t, err)
	})

	t.Run("DeleteWallet", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "wallets"`).
			WithArgs(walletID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(walletID)
		assert.NoError(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
