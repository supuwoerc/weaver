package utils

import (
	"testing"

	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/jwt"
	"github.com/supuwoerc/weaver/pkg/response"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAndTeardownTest(t *testing.T) func() {
	t.Log("Setting up test data...")
	return func() {
		t.Log("Tearing down test data...")
	}
}

func TestGetContextClaims(t *testing.T) {
	type args struct {
		ctx *gin.Context
	}
	emptyContext, _ := gin.CreateTestContext(nil)
	contextWithClaims, _ := gin.CreateTestContext(nil)
	claims := &jwt.TokenClaims{
		RegisteredClaims: jwt2.RegisteredClaims{},
		User: &jwt.TokenClaimsBasic{
			ID:       100,
			Email:    "test@email.com",
			Nickname: lo.ToPtr("test nickname"),
		},
	}
	contextWithClaims.Set(constant.ClaimsContextKey, claims)
	contextWithInvalidClaims, _ := gin.CreateTestContext(nil)
	contextWithInvalidClaims.Set(constant.ClaimsContextKey, nil)
	tests := []struct {
		name    string
		args    args
		want    *jwt.TokenClaims
		wantErr bool
		err     error
	}{
		{
			name:    "EmptyContext",
			args:    args{ctx: emptyContext},
			want:    nil,
			wantErr: true,
			err:     response.UserNotExist,
		},
		{
			name:    "ValidClaims",
			args:    args{ctx: contextWithClaims},
			want:    claims,
			wantErr: false,
			err:     nil,
		},
		{
			name:    "InvalidClaims",
			args:    args{ctx: contextWithInvalidClaims},
			want:    nil,
			wantErr: true,
			err:     response.UserNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setupAndTeardownTest(t)
			defer teardown()
			got, err := GetContextClaims(tt.args.ctx)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.err, err)
			assert.Equal(t, got, tt.want)
		})
	}
}
