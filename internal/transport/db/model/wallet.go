package model

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid"`
	Balance int64     `gorm:"not null;default:0"`
}
