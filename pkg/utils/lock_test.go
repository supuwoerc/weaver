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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
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
		err := r.lock.Lock(ctx, true)
		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
	})

	r.Run("lock with invalid state", func() {
		// 设置未锁定
		r.lock.state.Store(999)
		// 调用 Lock
		err := r.lock.Lock(ctx, true)
		// 验证结果
		assert.ErrorContains(t, err, "lock is in invalid state: 999")
		r.mutex.AssertNotCalled(t, "Lock")
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

func (r *RedisLockSuite) TestRedisLock_TryLock() {
	t := r.T()
	ctx := context.Background()

	r.Run("successful try lock", func() {
		// 设置未锁定
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - Lock() 成功
		r.mutex.On("TryLock").Return(nil)
		// 调用 Lock
		err := r.lock.TryLock(ctx, true)
		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
	})

	r.Run("try lock with invalid state", func() {
		// 设置未锁定
		r.lock.state.Store(999)
		// 调用 Lock
		err := r.lock.TryLock(ctx, true)
		// 验证结果
		assert.ErrorContains(t, err, "lock is in invalid state: 999")
		r.mutex.AssertNotCalled(t, "Lock")
	})

	r.Run("try lock already locked", func() {
		// 已经锁定
		r.lock.state.Store(lockStateLocked)
		// 调用 Lock
		err := r.lock.TryLock(ctx, false)
		// 验证结果
		assert.Equal(t, response.Busy, err)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		// 不应该调用底层的 Lock 方法
		r.mutex.AssertNotCalled(t, "Lock")
	})

	r.Run("try lock already released", func() {
		// 已经释放
		r.lock.state.Store(lockStateReleased)
		// 调用 Lock
		err := r.lock.TryLock(ctx, false)
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "lock has been released")
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		// 不应该调用底层的 Lock 方法
		r.mutex.AssertNotCalled(t, "Lock")
	})

	r.Run("try lock is extending", func() {
		// 正在续期
		r.lock.state.Store(lockStateExtending)
		// 调用 Lock
		err := r.lock.TryLock(ctx, false)
		// 验证结果
		assert.Equal(t, response.Busy, err)
		assert.Equal(t, int64(lockStateExtending), r.lock.state.Load())
		// 不应该调用底层的 Lock 方法
		r.mutex.AssertNotCalled(t, "Lock")
	})

	r.Run("try lock with redis lock returns ErrTaken", func() {
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - Lock() 返回 ErrTaken
		lockErr := &redsync.ErrTaken{}
		r.mutex.On("TryLock").Return(lockErr)
		// 调用 Lock
		err := r.lock.TryLock(ctx, false)
		// 验证结果
		assert.Equal(t, response.Busy, err)
		assert.Equal(t, int64(lockStateUnlocked), r.lock.state.Load()) // 状态应该恢复
		r.mutex.AssertExpectations(t)
	})

	r.Run("try lock with redis lock returns other error", func() {
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - Lock() 返回其他错误
		lockErr := errors.New("redis connection error")
		r.mutex.On("TryLock").Return(lockErr)
		// 调用 Lock
		err := r.lock.TryLock(ctx, false)
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis lock err")
		assert.Contains(t, err.Error(), "redis connection error")
		assert.Equal(t, int64(lockStateUnlocked), r.lock.state.Load()) // 状态应该恢复
		r.mutex.AssertExpectations(t)
	})

	r.Run("try lock concurrent lock attempts", func() {
		r.lock.state.Store(lockStateUnlocked)
		// 设置 mock 期望 - 第一次 Lock() 成功
		r.mutex.On("TryLock").Return(nil).Once()
		// 第一次调用应该成功
		err1 := r.lock.TryLock(ctx, false)
		assert.NoError(t, err1)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		// 第二次调用应该失败（锁已被占用）
		err2 := r.lock.TryLock(ctx, false)
		assert.Equal(t, response.Busy, err2)
		// 验证 mock 调用（应该只调用一次）
		r.mutex.AssertExpectations(t)
		r.mutex.AssertNumberOfCalls(t, "TryLock", 1)
	})

}

func (r *RedisLockSuite) TestRedisLock_Unlock() {
	t := r.T()

	r.Run("successful unlock from locked state", func() {
		// 设置锁定状态
		r.lock.state.Store(lockStateLocked)
		// 设置 mock 期望 - Unlock() 成功
		r.mutex.On("Unlock").Return(true, nil)
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		// 验证 stopChan 已关闭
		select {
		case <-r.lock.stopChan:
			// 期望的行为 - stopChan 应该被关闭
		default:
			t.Error("stopChan should be closed after unlock")
		}
		r.mutex.AssertExpectations(t)
	})

	r.Run("unlock from unlocked state", func() {
		// 设置未锁定状态
		r.lock.state.Store(lockStateUnlocked)
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unlock: lock not held")
		assert.Equal(t, int64(lockStateUnlocked), r.lock.state.Load())
		// 不应该调用底层的 Unlock 方法
		r.mutex.AssertNotCalled(t, "Unlock")
	})

	r.Run("unlock from already released state", func() {
		// 设置已释放状态
		r.lock.state.Store(lockStateReleased)
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		// 不应该调用底层的 Unlock 方法
		r.mutex.AssertNotCalled(t, "Unlock")
	})

	r.Run("unlock from extending state - extend completes in time", func() {
		// 设置正在扩展状态
		r.lock.state.Store(lockStateExtending)
		// 模拟 extend 操作完成
		go func() {
			time.Sleep(100 * time.Millisecond)
			r.lock.state.Store(lockStateLocked)
			close(r.lock.extendDone)
		}()
		// 设置 mock 期望 - Unlock() 成功
		r.mutex.On("Unlock").Return(true, nil)
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
	})

	r.Run("unlock from extending state - timeout waiting for extend", func() {
		// 设置正在扩展状态
		r.lock.state.Store(lockStateExtending)
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout waiting for extend to complete")
		assert.Equal(t, int64(lockStateExtending), r.lock.state.Load())
		// 不应该调用底层的 Unlock 方法
		r.mutex.AssertNotCalled(t, "Unlock")
	})

	r.Run("unlock with invalid state", func() {
		// 设置无效状态
		r.lock.state.Store(999)
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unlock: lock in invalid state: 999")
		assert.Equal(t, int64(999), r.lock.state.Load())
		r.mutex.AssertNotCalled(t, "Unlock")
	})

	r.Run("redis unlock returns error", func() {
		// 设置锁定状态
		r.lock.state.Store(lockStateLocked)
		// 设置 mock 期望 - Unlock() 返回错误
		unlockErr := errors.New("redis connection error")
		r.mutex.On("Unlock").Return(false, unlockErr)
		r.mutex.On("Name").Return("test-lock")
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis connection error")
		assert.Contains(t, err.Error(), "test-lock")
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
	})

	r.Run("redis unlock returns lock already expired error", func() {
		// 设置锁定状态
		r.lock.state.Store(lockStateLocked)
		// 设置 mock 期望 - Unlock() 返回锁已过期错误
		r.mutex.On("Unlock").Return(false, redsync.ErrLockAlreadyExpired)
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果 - 锁已过期时应该返回 nil（认为解锁成功）
		assert.NoError(t, err)
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
	})

	r.Run("redis unlock returns false", func() {
		// 设置锁定状态
		r.lock.state.Store(lockStateLocked)
		// 设置 mock 期望 - Unlock() 返回 false
		r.mutex.On("Unlock").Return(false, nil)
		r.mutex.On("Name").Return("test-lock")
		// 调用 Unlock
		err := r.lock.Unlock()
		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test-lock unlock failed")
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
	})

	r.Run("concurrent unlock attempts", func() {
		// 设置锁定状态
		r.lock.state.Store(lockStateLocked)
		// 设置 mock 期望 - 第一次 Unlock() 成功
		r.mutex.On("Unlock").Return(true, nil).Once()
		// 第一次调用应该成功
		err1 := r.lock.Unlock()
		assert.NoError(t, err1)
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		// 第二次调用应该直接返回 nil（已经释放）
		err2 := r.lock.Unlock()
		assert.NoError(t, err2)
		// 验证 mock 调用（应该只调用一次）
		r.mutex.AssertExpectations(t)
		r.mutex.AssertNumberOfCalls(t, "Unlock", 1)
		r.mutex.AssertNotCalled(t, "Name")
	})

}

func (r *RedisLockSuite) TestRedisLock_UnlockWithLog() {
	t := r.T()
	ctx := context.Background()

	r.Run("unlock succeeds - no log captured", func() {
		// 使用真实的测试 logger，但可以捕获日志输出
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 mock 期望 - Unlock() 成功
		r.mutex.On("Unlock").Return(true, nil)
		// 调用 unlockWithLog
		r.lock.unlockWithLog(ctx)
		// 验证 unlock 成功且状态正确
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		r.mutex.AssertExpectations(t)
		r.mutex.AssertNotCalled(t, "Name")
	})

	r.Run("unlock fails - error logged", func() {
		// 使用 observer 来捕获日志
		observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		observedLogger := zap.New(observedZapCore).Sugar()
		testLogger := logger.NewLogger(observedLogger)
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 unlock 失败
		unlockErr := errors.New("redis connection error")
		r.mutex.On("Unlock").Return(false, unlockErr)
		r.mutex.On("Name").Return("test-lock").Maybe()
		// 调用 unlockWithLog
		r.lock.unlockWithLog(ctx)
		// 验证状态仍然是 locked（因为 unlock 失败）
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		// 验证错误日志被记录
		logEntries := observedLogs.All()
		assert.Len(t, logEntries, 1)
		assert.Equal(t, "redis unlock fail", logEntries[0].Message)
		assert.Contains(t, logEntries[0].Context[0].String, "redis connection error")
		assert.Contains(t, logEntries[0].Context[0].String, "test-lock")
		r.mutex.AssertExpectations(t)
	})

	r.Run("unlock from unlocked state - error logged", func() {
		// 使用 observer 来捕获日志
		observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		observedLogger := zap.New(observedZapCore).Sugar()
		testLogger := logger.NewLogger(observedLogger)
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateUnlocked)
		// 调用 unlockWithLog
		r.lock.unlockWithLog(ctx)
		// 验证状态没有改变
		assert.Equal(t, int64(lockStateUnlocked), r.lock.state.Load())
		// 验证错误日志被记录
		logEntries := observedLogs.All()
		assert.Len(t, logEntries, 1)
		assert.Equal(t, "redis unlock fail", logEntries[0].Message)
		assert.Equal(t, "cannot unlock: lock not held", logEntries[0].Context[0].String)
		// 不应该调用底层的 Unlock 方法
		r.mutex.AssertNotCalled(t, "Unlock")
	})

	r.Run("unlock with redis returning lock already expired - no error logged", func() {
		observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		observedLogger := zap.New(observedZapCore).Sugar()
		testLogger := logger.NewLogger(observedLogger)
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 unlock 返回锁已过期错误
		r.mutex.On("Unlock").Return(false, redsync.ErrLockAlreadyExpired)
		// 调用 unlockWithLog
		r.lock.unlockWithLog(ctx)
		// 验证状态变为 released
		assert.Equal(t, int64(lockStateReleased), r.lock.state.Load())
		// 验证没有错误日志被记录
		logEntries := observedLogs.All()
		assert.Len(t, logEntries, 0)
		r.mutex.AssertExpectations(t)
	})
}

func (r *RedisLockSuite) TestRedisLock_Extend() {
	t := r.T()

	r.Run("stopChan is closed - should return false immediately", func() {
		ctx := context.Background()
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 关闭 stopChan
		close(r.lock.stopChan)
		// 调用 extend
		result := r.lock.extend(ctx)
		// 验证结果
		assert.False(t, result)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		// 不应该调用底层的 Extend 方法
		r.mutex.AssertNotCalled(t, "Extend")
	})

	r.Run("state is not lockStateLocked - should return false", func() {
		ctx := context.Background()
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())

		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}

		// 测试不同的非 lockStateLocked 状态
		testStates := []int64{lockStateUnlocked, lockStateReleased, lockStateExtending}
		for _, state := range testStates {
			r.lock.state.Store(state)
			// 调用 extend
			result := r.lock.extend(ctx)
			// 验证结果
			assert.False(t, result)
			assert.Equal(t, state, r.lock.state.Load()) // 状态应该不变
		}
		// 不应该调用底层的 Extend 方法
		r.mutex.AssertNotCalled(t, "Extend")
	})

	r.Run("context deadline approaching - should log warning and return false", func() {
		observedZapCore, observedLogs := observer.New(zapcore.WarnLevel)
		observedLogger := zap.New(observedZapCore).Sugar()
		testLogger := logger.NewLogger(observedLogger)
		// 创建一个即将到期的 context
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		defer cancel()
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 mock - Until() 返回一个更晚的时间
		futureTime := time.Now().Add(5 * time.Second)
		r.mutex.On("Until").Return(futureTime)
		r.mutex.On("Name").Return("test-lock")
		// 调用 extend
		result := r.lock.extend(ctx)
		// 验证结果
		assert.False(t, result)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load()) // 状态应该恢复
		// 验证警告日志被记录
		logEntries := observedLogs.All()
		assert.Len(t, logEntries, 1)
		assert.Equal(t, "skipping extend due to approaching deadline", logEntries[0].Message)
		assert.Equal(t, zapcore.WarnLevel, logEntries[0].Level)
		// 验证日志字段
		fields := logEntries[0].Context
		fieldMap := make(map[string]interface{})
		for _, field := range fields {
			fieldMap[field.Key] = field.Interface
		}
		assert.Contains(t, fieldMap, "name")
		assert.Contains(t, fieldMap, "deadline")
		assert.Contains(t, fieldMap, "mutex until")
		r.mutex.AssertExpectations(t)
	})

	r.Run("redis extend returns error - should log error and return false", func() {
		observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		observedLogger := zap.New(observedZapCore).Sugar()
		testLogger := logger.NewLogger(observedLogger)
		ctx := context.Background()
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 mock - Extend() 返回错误
		extendErr := errors.New("redis connection failed")
		r.mutex.On("Extend").Return(false, extendErr)
		r.mutex.On("Name").Return("test-lock")
		// 调用 extend
		result := r.lock.extend(ctx)
		// 验证结果
		assert.False(t, result)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load()) // 状态应该恢复
		// 验证错误日志被记录
		logEntries := observedLogs.All()
		assert.Len(t, logEntries, 1)
		assert.Equal(t, "lock couldn't be extended", logEntries[0].Message)
		assert.Equal(t, zapcore.ErrorLevel, logEntries[0].Level)
		// 验证日志字段
		fields := logEntries[0].Context
		assert.Equal(t, "redis connection failed", fields[0].String)
		assert.Equal(t, "test-lock", fields[1].String)
		r.mutex.AssertExpectations(t)
	})

	r.Run("redis extend returns false - should log error and return false", func() {
		observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		observedLogger := zap.New(observedZapCore).Sugar()
		testLogger := logger.NewLogger(observedLogger)
		ctx := context.Background()
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 mock - Extend() 返回 false
		r.mutex.On("Extend").Return(false, nil)
		r.mutex.On("Name").Return("test-lock")
		// 调用 extend
		result := r.lock.extend(ctx)
		// 验证结果
		assert.False(t, result)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load()) // 状态应该恢复
		// 验证错误日志被记录
		logEntries := observedLogs.All()
		assert.Len(t, logEntries, 1)
		assert.Equal(t, "redis extend lock fail ", logEntries[0].Message)
		assert.Equal(t, zapcore.ErrorLevel, logEntries[0].Level)
		// 验证日志字段
		fields := logEntries[0].Context
		assert.Equal(t, "test-lock", fields[0].String)
		r.mutex.AssertExpectations(t)
	})

	r.Run("successful extend - should return true", func() {
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		ctx := context.Background()
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 mock - Extend() 成功
		r.mutex.On("Extend").Return(true, nil)
		// 调用 extend
		result := r.lock.extend(ctx)
		// 验证结果
		assert.True(t, result)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load()) // 状态应该恢复到 lockStateLocked
		r.mutex.AssertExpectations(t)
	})

	r.Run("state transitions correctly during extend", func() {
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		ctx := context.Background()
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 用于跟踪状态变化的 channel
		stateChanges := make(chan int64, 1)
		defer close(stateChanges)
		// 设置一个会阻塞的 mock，让我们能观察状态变化
		r.mutex.On("Extend").Run(func(args mock.Arguments) {
			// 在 Extend 调用期间，状态应该是 lockStateExtending
			stateChanges <- r.lock.state.Load()
		}).Return(true, nil)
		// 在另一个 goroutine 中启动 extend
		resultChan := make(chan bool)
		defer close(resultChan)
		go func() {
			result := r.lock.extend(ctx)
			resultChan <- result
		}()
		// 等待 extend 完成
		result := <-resultChan
		// 验证结果
		assert.True(t, result)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load()) // 最终状态应该恢复
		// 验证在 Extend 期间状态是 lockStateExtending
		select {
		case state := <-stateChanges:
			assert.Equal(t, int64(lockStateExtending), state)
		default:
			t.Error("Should have captured state change during extend")
		}
		r.mutex.AssertExpectations(t)
	})

	r.Run("extendDone channel is managed correctly", func() {
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		ctx := context.Background()
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 mock
		r.mutex.On("Extend").Return(true, nil)
		// 验证 extendDone 初始状态（应该是开放的）
		select {
		case <-r.lock.extendDone:
			t.Error("extendDone should not be closed initially")
		default:
			// 期望的行为
		}
		// 调用 extend
		result := r.lock.extend(ctx)
		// 验证结果
		assert.True(t, result)
		// 验证 extendDone 被关闭了（defer 函数应该关闭它）
		select {
		case <-r.lock.extendDone:
			// 期望的行为 - channel 应该被关闭
		default:
			t.Error("extendDone should be closed after extend completes")
		}
		r.mutex.AssertExpectations(t)
		r.mutex.AssertNumberOfCalls(t, "Extend", 1)
	})

	r.Run("context without deadline - should proceed normally", func() {
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		ctx := context.Background() // 没有 deadline 的 context
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      defaultLockDuration,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 mock - Extend() 成功
		r.mutex.On("Extend").Return(true, nil)
		// 调用 extend
		result := r.lock.extend(ctx)
		// 验证结果
		assert.True(t, result)
		assert.Equal(t, int64(lockStateLocked), r.lock.state.Load())
		// 不应该调用 Until() 因为没有 deadline
		r.mutex.AssertNotCalled(t, "Until")
		r.mutex.AssertExpectations(t)
	})
}

// ... existing code ...

func (r *RedisLockSuite) TestRedisLock_AutoExtend() {
	t := r.T()

	r.Run("context canceled immediately - should call unlockWithLog", func() {
		testLogger := logger.NewLogger(zaptest.NewLogger(t).Sugar())
		// 创建一个已经取消的 context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      100 * time.Millisecond, // 使用较短的时间便于测试
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 unlock mock
		r.mutex.On("Unlock").Return(true, nil)
		// 用于等待 autoExtend 完成的 channel
		done := make(chan struct{})
		defer close(done)
		// 启动 autoExtend
		go func() {
			defer func() { done <- struct{}{} }()
			r.lock.autoExtend(ctx)
		}()
		// 等待 autoExtend 完成
		select {
		case <-done:
			// 期望的行为
		case <-time.After(time.Second):
			t.Error("autoExtend should complete quickly when context is canceled")
		}
		// 验证 unlock 被调用
		r.mutex.AssertExpectations(t)
		r.mutex.AssertNumberOfCalls(t, "Unlock", 1)
	})

	r.Run("stopChan closed - should exit without unlocking", func() {
		ctx := context.Background()
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      100 * time.Millisecond,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 关闭 stopChan
		close(r.lock.stopChan)
		// 用于等待 autoExtend 完成的 channel
		done := make(chan struct{})
		defer close(done)
		// 启动 autoExtend
		go func() {
			defer func() { done <- struct{}{} }()
			r.lock.autoExtend(ctx)
		}()

		// 等待 autoExtend 完成
		select {
		case <-done:
			// 期望的行为
		case <-time.After(1 * time.Second):
			t.Error("autoExtend should complete quickly when stopChan is closed")
		}

		// 不应该调用 unlock
		r.mutex.AssertNotCalled(t, "Unlock")
	})

	// FIXME test case
	r.Run("ticker triggers extend - successful case", func() {
		ctx := context.Background()
		testLogger := logger.NewLogger(zaptest.NewLogger(r.T()).Sugar())
		// 重新初始化 lock，使用非常短的 duration 使 ticker 快速触发
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      1 * time.Second,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 extend mock - 第一次成功，第二次失败触发退出
		r.mutex.On("Extend").Return(true, nil).Once()  // 第一次成功
		r.mutex.On("Extend").Return(false, nil).Once() // 第二次失败
		r.mutex.On("Name").Return("test-lock")
		// 用于等待 autoExtend 完成的 channel
		done := make(chan struct{})
		defer close(done)
		// 启动 autoExtend
		go func() {
			defer func() { done <- struct{}{} }()
			r.lock.autoExtend(ctx)
		}()
		// 等待 autoExtend 完成（应该在第二次 extend 失败后退出）
		select {
		case <-done:
			// 期望的行为
		case <-time.After(2 * time.Second):
			t.Error("autoExtend should complete after extend fails")
		}
		// 验证 extend 被调用了两次
		r.mutex.AssertExpectations(t)
		r.mutex.AssertNumberOfCalls(t, "Extend", 2)
		r.mutex.AssertNumberOfCalls(t, "Name", 1)
	})

	r.Run("extend fails - should exit autoExtend", func() {
		observedZapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		observedLogger := zap.New(observedZapCore).Sugar()
		testLogger := logger.NewLogger(observedLogger)

		ctx := context.Background()

		// 重新初始化 lock
		r.lock = &RedisLock{
			Mutex:         r.mutex,
			duration:      30 * time.Millisecond,
			logger:        testLogger,
			stopChan:      make(chan struct{}),
			extendDone:    make(chan struct{}),
			extendTimeout: defaultExtendTimeout,
		}
		r.lock.state.Store(lockStateLocked)
		// 设置 extend mock - 失败
		r.mutex.On("Extend").Return(false, errors.New("extend failed"))
		r.mutex.On("Name").Return("test-lock")
		// 用于等待 autoExtend 完成的 channel
		done := make(chan struct{})
		defer close(done)
		// 启动 autoExtend
		go func() {
			defer func() { done <- struct{}{} }()
			r.lock.autoExtend(ctx)
		}()
		// 等待 autoExtend 完成
		select {
		case <-done:
			// 期望的行为
		case <-time.After(1 * time.Second):
			t.Error("autoExtend should complete when extend fails")
		}
		// 验证 extend 被调用并且有错误日志
		r.mutex.AssertExpectations(t)
		r.mutex.AssertCalled(t, "Extend")
		// 验证错误日志
		logEntries := observedLogs.All()
		assert.Equal(t, len(logEntries), 1)
		assert.Contains(t, logEntries[0].Entry.Message, "lock couldn't be extended")
	})
}
