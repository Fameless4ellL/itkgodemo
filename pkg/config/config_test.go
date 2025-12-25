package config

import (
	"testing"
)

func TestInit(t *testing.T) {

	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("PORT", "8080")
	t.Setenv("DEBUG", "true")

	Init()

	if Port != 8080 {
		t.Errorf("expected Port 8080, got %d", Port)
	}
	if Debug != true {
		t.Errorf("expected Debug true, got %v", Debug)
	}
}
