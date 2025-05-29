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
		assert.PanicsWithValue(t, panicMessage, func() {
			_ = manager.Refresh(ctx, "cache3")
		})
	})
}

func Test_operateCache(t *testing.T) {
	ctx := context.Background()
	cache1 := &testCache{}
	cache1.On("Key").Return("cache1")
	cache1.On("Clean", ctx).Return(nil)
	cache2 := &testCache{}
	cache2.On("Key").Return("cache2")
	cache2.On("Clean", ctx).Return(fmt.Errorf("cache2:test clean err"))
	cache3 := &testCache{}
	cache3.On("Key").Return("cache3")
	cache3.On("Clean", ctx).Panic("cache3:test clean panic")
	caches := []SystemCache{cache1, cache2, cache3}
	manager := NewSystemCacheManager(caches...)
	t.Run("operateCache with invalid op", func(t *testing.T) {
		err := operateCache(ctx, cacheOperate(100), manager, cache1.Key())
		assert.ErrorContains(t, err, "is invalid operate")
	})
	t.Run("operateCache with clean success op", func(t *testing.T) {
		err := operateCache(ctx, clean, manager, cache1.Key())
		assert.NoError(t, err)
	})
	t.Run("operateCache with clean fail op", func(t *testing.T) {
		err := operateCache(ctx, clean, manager, cache2.Key())
		assert.ErrorContains(t, err, "cache2:test clean err")
	})
	t.Run("operateCache with clean panic op", func(t *testing.T) {
		assert.PanicsWithValue(t, "cache3:test clean panic", func() {
			_ = operateCache(ctx, clean, manager, cache3.Key())
		})
	})
}
