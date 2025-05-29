package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 定义环境常量，避免魔法字符串
const (
	EnvProd = "prod"
	EnvDev  = "dev"
	EnvTest = "test"
)

func TestConfig_Environment(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name   string
		env    string
		isProd bool
		isDev  bool
		isTest bool
	}{
		{
			name:   "production environment",
			env:    EnvProd,
			isProd: true,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "development environment",
			env:    EnvDev,
			isProd: false,
			isDev:  true,
			isTest: false,
		},
		{
			name:   "test environment",
			env:    EnvTest,
			isProd: false,
			isDev:  false,
			isTest: true,
		},
		{
			name:   "empty environment",
			env:    "",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "invalid environment",
			env:    "staging",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "case sensitive - PROD",
			env:    "PROD",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "case sensitive - DEV",
			env:    "DEV",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "case sensitive - TEST",
			env:    "TEST",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "whitespace - prod with spaces",
			env:    " prod ",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "whitespace - dev with spaces",
			env:    " dev ",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
		{
			name:   "whitespace - test with spaces",
			env:    " test ",
			isProd: false,
			isDev:  false,
			isTest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Env: tt.env,
			}

			// 测试 IsProd
			got := c.IsProd()
			require.Equal(t, tt.isProd, got,
				"IsProd() = %v, want %v for env=%q", got, tt.isProd, tt.env)

			// 测试 IsDev
			got = c.IsDev()
			require.Equal(t, tt.isDev, got,
				"IsDev() = %v, want %v for env=%q", got, tt.isDev, tt.env)

			// 测试 IsTest
			got = c.IsTest()
			require.Equal(t, tt.isTest, got,
				"IsTest() = %v, want %v for env=%q", got, tt.isTest, tt.env)

			// 验证环境互斥性
			envChecks := 0
			if c.IsProd() {
				envChecks++
			}
			if c.IsDev() {
				envChecks++
			}
			if c.IsTest() {
				envChecks++
			}
			assert.LessOrEqual(t, envChecks, 1,
				"Environment checks should be mutually exclusive for env=%q", tt.env)
		})
	}
}
