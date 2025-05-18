package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type testCache struct {
	mock.Mock
}

func (t *testCache) Key() string {
	called := t.Called()
	return called.String(0)
}

func (t *testCache) Refresh(ctx context.Context) error {
	called := t.Called(ctx)
	return called.Error(0)
}

func (t *testCache) Clean(ctx context.Context) error {
	called := t.Called(ctx)
	return called.Error(0)
}

func TestNewSystemCacheManager(t *testing.T) {
	c := &testCache{}
	type args struct {
		caches []SystemCache
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				caches: []SystemCache{c},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSystemCacheManager(tt.args.caches...)
			assert.Equal(t, len(got.caches), tt.want)
			assert.Equal(t, got.caches[0], c)
		})
	}
}

func TestSystemCacheManager_Refresh(t *testing.T) {
	cache1 := &testCache{} // mock缓存1
	cache2 := &testCache{} // mock缓存2
	cache3 := &testCache{} // mock缓存3
	ctx := context.Background()

	cache1.On("Key").Return("cache1")
	cache2.On("Key").Return("cache2")
	cache3.On("Key").Return("cache3")

	cache1.On("Refresh", ctx).Return(nil)
	refreshErr := fmt.Errorf("cache2:refresh err")
	cache2.On("Refresh", ctx).Return(refreshErr)
	panicMessage := "cache3:refresh panic"
	cache3.On("Refresh", ctx).Panic(panicMessage)
	t.Run("RefreshWithoutTarget", func(t *testing.T) {
		caches := []SystemCache{cache1}
		manager := NewSystemCacheManager(caches...)
		err := manager.Refresh(ctx)
		assert.NoError(t, err)
	})
	t.Run("RefreshEmptyCache", func(t *testing.T) {
		manager := NewSystemCacheManager()
		err := manager.Refresh(ctx, "cache1")
		assert.ErrorContains(t, err, "empty")
	})
	t.Run("RefreshNonExistCache", func(t *testing.T) {
		caches := []SystemCache{cache1}
		manager := NewSystemCacheManager(caches...)
		err := manager.Refresh(ctx, "cache3")
		assert.ErrorContains(t, err, "refresh cache fail: cache cache3 not found")
	})
	t.Run("RefreshExistCache", func(t *testing.T) {
		caches := []SystemCache{cache1}
		manager := NewSystemCacheManager(caches...)
		err := manager.Refresh(ctx, "cache1")
		assert.NoError(t, err)
	})
	t.Run("RefreshCacheWithErr", func(t *testing.T) {
		caches := []SystemCache{cache2}
		manager := NewSystemCacheManager(caches...)
		err := manager.Refresh(ctx, "cache2")
		assert.ErrorIs(t, err, refreshErr)
	})
	t.Run("RefreshCacheWithPanic", func(t *testing.T) {
		caches := []SystemCache{cache3}
		manager := NewSystemCacheManager(caches...)
		assert.Panics(t, func() {
			_ = manager.Refresh(ctx, "cache3")
		})
	})
}
