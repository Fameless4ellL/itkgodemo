package config

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	Init()

	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("PORT", "8080")
	t.Setenv("DEBUG", "true")

	if Port != 0 {
		t.Errorf("expected Port 8080, got %d", Port)
	}
	if Debug != false {
		t.Errorf("expected Debug true, got %v", Debug)
	}
}

func TestInitDefaults(t *testing.T) {
	Init()

	os.Unsetenv("PORT")
	os.Unsetenv("DEBUG")

	if Port != 0 {
		t.Errorf("expected default Port 0, got %d", Port)
	}
	if Debug != false {
		t.Errorf("expected default Debug false, got %v", Debug)
	}
}
