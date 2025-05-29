package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/supuwoerc/weaver/pkg/constant"
	"gorm.io/gorm"
)

func TestLoadManager(t *testing.T) {
	manager := &TransactionManager{
		DB:                           &gorm.DB{},
		AlreadyCommittedOrRolledBack: false,
	}
	tests := []struct {
		name    string
		ctx     context.Context
		want    *TransactionManager
		wantNil bool
	}{
		{
			name: "valid transaction manager",
			ctx: context.WithValue(
				context.Background(),
				constant.TransactionKey,
				manager,
			),
			want:    manager,
			wantNil: false,
		},
		{
			name:    "nil context",
			ctx:     context.Background(),
			want:    nil,
			wantNil: true,
		},
		{
			name: "wrong type in context",
			ctx: context.WithValue(
				context.Background(),
				constant.TransactionKey,
				"not a transaction manager",
			),
			want:    nil,
			wantNil: true,
		},
		{
			name: "nil value in context",
			ctx: context.WithValue(
				context.Background(),
				constant.TransactionKey,
				nil,
			),
			want:    nil,
			wantNil: true,
		},
		{
			name: "different key in context",
			ctx: context.WithValue(
				context.Background(),
				constant.ContextKey("invalid_key"),
				&TransactionManager{},
			),
			want:    nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadManager(tt.ctx)
			if tt.wantNil {
				require.Nil(t, got, "LoadManager() should return nil")
				return
			}
			require.NotNil(t, got, "LoadManager() should not return nil")
			if tt.want.DB == nil {
				require.Nil(t, got.DB)
			} else {
				require.NotNil(t, got.DB)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoadManager_WithInject(t *testing.T) {
	t.Run("inject and load", func(t *testing.T) {
		manager := &TransactionManager{
			DB:                           &gorm.DB{},
			AlreadyCommittedOrRolledBack: false,
		}
		ctx := InjectManager(context.Background(), manager)
		loaded := LoadManager(ctx)
		assert.NotNil(t, loaded)
		assert.Equal(t, manager, loaded)
	})

	t.Run("inject nil manager", func(t *testing.T) {
		ctx := InjectManager(context.Background(), nil)
		loaded := LoadManager(ctx)
		assert.Nil(t, loaded)
	})

	t.Run("multiple inject operations", func(t *testing.T) {
		manager1 := &TransactionManager{
			DB:                           &gorm.DB{},
			AlreadyCommittedOrRolledBack: false,
		}
		manager2 := &TransactionManager{
			DB:                           &gorm.DB{},
			AlreadyCommittedOrRolledBack: true,
		}
		ctx := InjectManager(context.Background(), manager1)
		loaded1 := LoadManager(ctx)
		assert.Equal(t, manager1, loaded1)
		ctx = InjectManager(ctx, manager2)
		loaded2 := LoadManager(ctx)
		assert.Equal(t, manager2, loaded2)
	})
}

func TestFuzzKeyword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "normal string",
			input:    "test",
			expected: "%test%",
		},
		{
			name:     "string with %",
			input:    "test%query",
			expected: "%test\\%query%",
		},
		{
			name:     "string with _",
			input:    "test_query",
			expected: "%test\\_query%",
		},
		{
			name:     "string with both % and _",
			input:    "test%query_string",
			expected: "%test\\%query\\_string%",
		},
		{
			name:     "multiple % characters",
			input:    "%%%",
			expected: "%\\%\\%\\%%",
		},
		{
			name:     "multiple _ characters",
			input:    "___",
			expected: "%\\_\\_\\_%",
		},
		{
			name:     "mixed special characters",
			input:    "%_test%_",
			expected: "%\\%\\_test\\%\\_%",
		},
		{
			name:     "chinese characters",
			input:    "测试",
			expected: "%测试%",
		},
		{
			name:     "chinese with special chars",
			input:    "测试%结果_",
			expected: "%测试\\%结果\\_%",
		},
		{
			name:     "numbers",
			input:    "123",
			expected: "%123%",
		},
		{
			name:     "special characters only",
			input:    "%_",
			expected: "%\\%\\_%",
		},
		{
			name:     "spaces",
			input:    "test query",
			expected: "%test query%",
		},
		{
			name:     "spaces with special chars",
			input:    "test % query _ test",
			expected: "%test \\% query \\_ test%",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "%a%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FuzzKeyword(tt.input)
			require.Equal(t, tt.expected, result,
				"FuzzKeyword(%q) = %q, want %q",
				tt.input, result, tt.expected)
		})
	}
}
