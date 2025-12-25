package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name          string
		debugEnv      string
		expectedLevel logrus.Level
		expectFullTS  bool
	}{
		{
			name:          "Debug mode enabled",
			debugEnv:      "true",
			expectedLevel: logrus.DebugLevel,
			expectFullTS:  true,
		},
		{
			name:          "Debug mode disabled",
			debugEnv:      "false",
			expectedLevel: logrus.InfoLevel,
			expectFullTS:  false,
		},
		{
			name:          "Invalid debug value defaults to info",
			debugEnv:      "not_a_boolean",
			expectedLevel: logrus.InfoLevel,
			expectFullTS:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменную окружения для конкретного подтеста
			t.Setenv("DEBUG", tt.debugEnv)

			// Вызываем инициализацию
			Init()

			// Проверяем уровень логирования
			assert.Equal(t, tt.expectedLevel, Log.GetLevel())

			// Проверяем настройки форматтера через приведение типов
			formatter, ok := Log.Formatter.(*logrus.TextFormatter)
			assert.True(t, ok, "Formatter should be of type *logrus.TextFormatter")

			if tt.expectFullTS {
				assert.True(t, formatter.FullTimestamp)
			} else {
				assert.True(t, formatter.DisableTimestamp)
			}
		})
	}
}
