package utils

import (
	"testing"

	"github.com/go-redsync/redsync/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/redis"
	"go.uber.org/zap/zaptest"
)

func TestNewRedisLocksmith(t *testing.T) {
	t.Run("RedisLocksmith fields", func(t *testing.T) {
		l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		client := &redis.CommonRedisClient{}
		locksmith := NewRedisLocksmith(l, client)
		assert.NotNil(t, locksmith)
		assert.Equal(t, locksmith.logger, l)
		assert.Equal(t, locksmith.redisClient, client)
	})
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	called := m.Called(name, options)
	return called.Get(0).(*redsync.Mutex)
}

func TestRedisLocksmith_NewLock(t *testing.T) {
	l := logger.NewLogger(zaptest.NewLogger(t).Sugar())
	mockClient := &MockRedisClient{}
	locksmith := &RedisLocksmith{
		logger:      l,
		redisClient: mockClient,
	}
	t.Run("NewLock with single object", func(t *testing.T) {
		// 期望的锁名称
		expectedName := "lock:permission:id:100"
		mockMutex := &redsync.Mutex{}
		// 设置 mock 期望
		mockClient.On("NewMutex", expectedName, mock.AnythingOfType("[]redsync.Option")).Return(mockMutex)
		// 调用 NewLock
		lock := locksmith.NewLock(constant.PermissionIdPrefix, "100")
		// 验证返回的锁
		assert.NotNil(t, lock)
		assert.Equal(t, mockMutex, lock.Mutex)
		assert.Equal(t, defaultLockDuration, lock.duration)
		assert.Equal(t, l, lock.logger)
		assert.Equal(t, defaultExtendTimeout, lock.extendTimeout)
		assert.NotNil(t, lock.stopChan)
		assert.NotNil(t, lock.extendDone)
		assert.Equal(t, int64(lockStateUnlocked), lock.state.Load())
		// 验证 mock 调用
		mockClient.AssertExpectations(t)
	})

	t.Run("NewLock with multiple objects", func(t *testing.T) {
		// 期望的锁名称
		expectedName := "lock:permission:id:100_200_300"
		mockMutex := &redsync.Mutex{}
		// 设置 mock 期望
		mockClient.On("NewMutex", expectedName, mock.AnythingOfType("[]redsync.Option")).Return(mockMutex)
		// 调用 NewLock
		lock := locksmith.NewLock(constant.PermissionIdPrefix, "100", "200", "300")
		// 验证返回的锁
		assert.NotNil(t, lock)
		assert.Equal(t, mockMutex, lock.Mutex)
		// 验证 mock 调用
		mockClient.AssertExpectations(t)
	})

	t.Run("NewLock with different prefix types", func(t *testing.T) {
		testCases := []struct {
			name         string
			prefix       constant.Prefix
			objects      []string
			expectedName string
		}{
			{
				name:         "permission id prefix",
				prefix:       constant.PermissionIdPrefix,
				objects:      []string{"123"},
				expectedName: "lock:permission:id:123",
			},
			{
				name:         "role name prefix",
				prefix:       constant.RoleNamePrefix,
				objects:      []string{"admin"},
				expectedName: "lock:role:name:admin",
			},
			{
				name:         "user id prefix",
				prefix:       constant.UserIdPrefix,
				objects:      []string{"456"},
				expectedName: "lock:user:id:456",
			},
			{
				name:         "department name prefix",
				prefix:       constant.DepartmentNamePrefix,
				objects:      []string{"IT"},
				expectedName: "lock:department:name:IT",
			},
			{
				name:         "signup email prefix",
				prefix:       constant.SignUpEmailPrefix,
				objects:      []string{"test@example.com"},
				expectedName: "lock:signup:email:test@example.com",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockMutex := &redsync.Mutex{}
				mockClient.On("NewMutex", tc.expectedName, mock.AnythingOfType("[]redsync.Option")).Return(mockMutex)
				lock := locksmith.NewLock(tc.prefix, tc.objects...)
				assert.NotNil(t, lock)
				assert.Equal(t, mockMutex, lock.Mutex)
				mockClient.AssertExpectations(t)
			})
		}
	})

	t.Run("verify default values are set correctly", func(t *testing.T) {
		mockMutex := &redsync.Mutex{}
		mockClient.On("NewMutex", mock.AnythingOfType("string"), mock.AnythingOfType("[]redsync.Option")).Return(mockMutex)
		lock := locksmith.NewLock(constant.PermissionIdPrefix, "test")
		// 验证所有默认值
		assert.Equal(t, defaultLockDuration, lock.duration)
		assert.Equal(t, defaultExtendTimeout, lock.extendTimeout)
		assert.Equal(t, int64(lockStateUnlocked), lock.state.Load())
		// 验证 channel 已初始化
		assert.NotNil(t, lock.stopChan)
		assert.NotNil(t, lock.extendDone)
		// 验证可以从 channel 读取（应该是非阻塞的关闭状态）
		select {
		case <-lock.stopChan:
			t.Error("stopChan should not be closed initially")
		default:
		}
		mockClient.AssertExpectations(t)
	})

	t.Run("multiple NewLock calls create independent locks", func(t *testing.T) {
		mockMutex1 := &redsync.Mutex{}
		mockMutex2 := &redsync.Mutex{}
		mockClient.On("NewMutex", "lock:permission:id:100", mock.AnythingOfType("[]redsync.Option")).Return(mockMutex1)
		mockClient.On("NewMutex", "lock:role:id:200", mock.AnythingOfType("[]redsync.Option")).Return(mockMutex2)
		lock1 := locksmith.NewLock(constant.PermissionIdPrefix, "100")
		lock2 := locksmith.NewLock(constant.RoleIdPrefix, "200")
		// 验证两个锁是独立的
		assert.NotEqual(t, lock1, lock2)
		assert.NotSame(t, lock1.Mutex, lock2.Mutex)
		assert.NotEqual(t, lock1.stopChan, lock2.stopChan)
		assert.NotEqual(t, lock1.extendDone, lock2.extendDone)
		// 共享相同的 logger
		assert.Same(t, lock1.logger, lock2.logger)
		mockClient.AssertExpectations(t)
	})

}
