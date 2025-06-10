package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/redis"
	"github.com/supuwoerc/weaver/pkg/response"
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

type MockMutex struct {
	mock.Mock
}

func (m *MockMutex) Lock() error {
	called := m.Called()
	return called.Error(0)
}

func (m *MockMutex) TryLock() error {
	called := m.Called()
	return called.Error(0)
}

func (m *MockMutex) Unlock() (bool, error) {
	called := m.Called()
	return called.Bool(0), called.Error(1)
}

func (m *MockMutex) Extend() (bool, error) {
	called := m.Called()
	return called.Bool(0), called.Error(1)
}

func (m *MockMutex) Name() string {
	called := m.Called()
	return called.String(0)
}

func (m *MockMutex) Until() time.Time {
	called := m.Called()
	return called.Get(0).(time.Time)
}

type RedisLockSuite struct {
	suite.Suite
	mutex *MockMutex
	lock  *RedisLock
}

func TestRedisLockSuite(t *testing.T) {
	suite.Run(t, new(RedisLockSuite))
}

func (r *RedisLockSuite) SetupSubTest() {
	mockMutex := &MockMutex{}
	testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
	lock := &RedisLock{
		Mutex:         mockMutex,
		duration:      defaultLockDuration,
		logger:        testLogger,
		stopChan:      make(chan struct{}),
		extendDone:    make(chan struct{}),
		extendTimeout: defaultExtendTimeout,
	}
	r.mutex = mockMutex
	r.lock = lock
}

func (r *RedisLockSuite) TestRedisLock_Lock() {
	t := r.T()
	ctx := context.Background()

	r.Run("successful lock", func() {
		// 设置未锁定
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - Lock() 成功
		r.mutex.On("Lock").Return(nil)
		// 调用 Lock
		err := r.lock.Lock(ctx, false)
		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
	})

	r.Run("lock already locked", func() {
		// 已经锁定
		r.lock.state.Store(lockStateLocked)
		// 调用 Lock
		err := r.lock.Lock(ctx, false)
		// 验证结果
		assert.Equal(t, response.Busy, err)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		// 不应该调用底层的 Lock 方法
		r.mutex.AssertNotCalled(t, "Lock")
	})

	r.Run("lock already released", func() {
		// 已经释放
		r.lock.state.Store(lockStateReleased)
		// 调用 Lock
		err := r.lock.Lock(ctx, false)
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "lock has been released")
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		// 不应该调用底层的 Lock 方法
		r.mutex.AssertNotCalled(t, "Lock")
	})

	r.Run("lock is extending", func() {
		// 正在续期
		r.lock.state.Store(lockStateExtending)
		// 调用 Lock
		err := r.lock.Lock(ctx, false)
		// 验证结果
		assert.Equal(t, response.Busy, err)
		assert.Equal(t, int64(lockStateExtending), r.lock.state.Load())
		// 不应该调用底层的 Lock 方法
		r.mutex.AssertNotCalled(t, "Lock")
	})

	r.Run("redis lock returns ErrTaken", func() {
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - Lock() 返回 ErrTaken
		lockErr := &redsync.ErrTaken{}
		r.mutex.On("Lock").Return(lockErr)
		// 调用 Lock
		err := r.lock.Lock(ctx, false)
		// 验证结果
		assert.Equal(t, response.Busy, err)
		assert.Equal(t, int64(lockStateUnlocked), r.lock.state.Load()) // 状态应该恢复
		r.mutex.AssertExpectations(t)
	})

	r.Run("redis lock returns other error", func() {
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - Lock() 返回其他错误
		lockErr := errors.New("redis connection error")
		r.mutex.On("Lock").Return(lockErr)
		// 调用 Lock
		err := r.lock.Lock(ctx, false)
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis lock err")
		assert.Contains(t, err.Error(), "redis connection error")
		assert.Equal(t, int64(lockStateUnlocked), r.lock.state.Load()) // 状态应该恢复
		r.mutex.AssertExpectations(t)
	})

	r.Run("concurrent lock attempts", func() {
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - 第一次 Lock() 成功
		r.mutex.On("Lock").Return(nil).Once()
		// 第一次调用应该成功
		err1 := r.lock.Lock(ctx, false)
		assert.NoError(t, err1)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		// 第二次调用应该失败（锁已被占用）
		err2 := r.lock.Lock(ctx, false)
		assert.Equal(t, response.Busy, err2)
		// 验证 mock 调用（应该只调用一次）
		r.mutex.AssertExpectations(t)
		r.mutex.AssertNumberOfCalls(t, "Lock", 1)
	})

}
