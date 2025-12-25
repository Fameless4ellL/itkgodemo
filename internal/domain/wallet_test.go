package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()
	if wallet == nil {
		t.Fatal("Expected wallet pointer, got nil")
	}
	expectedBalance := int64(1000)
	if wallet.Balance != expectedBalance {
		t.Errorf("Expected initial balance %d, but got %d", expectedBalance, wallet.Balance)
	}
	if wallet.ID == uuid.Nil {
		t.Error("Expected generated UUID, but got uuid.Nil")
	}
}
