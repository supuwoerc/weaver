package job

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
	"go.uber.org/zap/zaptest"
)

// MockSystemJob 用于测试的模拟任务
type MockSystemJob struct {
	name           string
	ifStillRunning constant.JobStillMode
	interval       string
	handleCalled   bool
	handleCount    int
}

func (m *MockSystemJob) Name() string {
	return m.name
}

func (m *MockSystemJob) IfStillRunning() constant.JobStillMode {
	return m.ifStillRunning
}

func (m *MockSystemJob) Handle() {
	m.handleCalled = true
}

func (m *MockSystemJob) Interval() string {
	return m.interval
}

func TestNewSystemJobManager(t *testing.T) {
	zapLogger := logger.NewLogger(zaptest.NewLogger(t).Sugar())
	emailClient := &initialize.EmailClient{}
	cronLogger := initialize.NewCronLogger(zapLogger, emailClient)
	cronClient := initialize.NewCronClient(cronLogger)
	t.Parallel()
	t.Run("empty jobs SystemJobManager", func(t *testing.T) {
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger)
		require.NotNil(t, manager)
		assert.Equal(t, cronLogger, manager.cronLogger)
		assert.Equal(t, cronClient, manager.cronClient)
		assert.Equal(t, zapLogger, manager.logger)
		assert.NotNil(t, manager.jobsMap)
		assert.Empty(t, manager.jobsMap)
		assert.Empty(t, manager.jobs)
	})

	t.Run("SystemJobManager with jobs", func(t *testing.T) {
		mockJob1 := &MockSystemJob{
			name:           "test-job-1",
			ifStillRunning: constant.Skip,
			interval:       "*/5 * * * * *",
		}

		mockJob2 := &MockSystemJob{
			name:           "test-job-2",
			ifStillRunning: constant.Delay,
			interval:       "0 */1 * * * *",
		}

		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob1, mockJob2)

		require.NotNil(t, manager)
		assert.Equal(t, cronLogger, manager.cronLogger)
		assert.Equal(t, cronClient, manager.cronClient)
		assert.Equal(t, zapLogger, manager.logger)
		assert.NotNil(t, manager.jobsMap)
		assert.Empty(t, manager.jobsMap) // jobsMap 在注册任务时才会填充
		assert.Len(t, manager.jobs, 2)
		assert.Contains(t, manager.jobs, mockJob1)
		assert.Contains(t, manager.jobs, mockJob2)
	})

	t.Run("SystemJobManager contains none mode job", func(t *testing.T) {
		mockJob := &MockSystemJob{
			name:           "test-job-none",
			ifStillRunning: constant.None,
			interval:       "0 0 * * * *",
		}
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob)
		require.NotNil(t, manager)
		assert.Len(t, manager.jobs, 1)
		assert.Contains(t, manager.jobs, mockJob)
	})
}

func TestSystemJobManager_skip(t *testing.T) {
	zapLogger := logger.NewLogger(zaptest.NewLogger(t).Sugar())
	emailClient := &initialize.EmailClient{}
	cronLogger := initialize.NewCronLogger(zapLogger, emailClient)
	cronClient := initialize.NewCronClient(cronLogger)
	manager := NewSystemJobManager(cronLogger, cronClient, zapLogger)

	t.Run("skip wrap func", func(t *testing.T) {
		executed := false
		testFunc := func() {
			executed = true
		}
		job := manager.skip(testFunc)
		require.NotNil(t, job)
		// 执行任务
		job.Run()
		assert.True(t, executed)
	})

	t.Run("skip func executed multiple times", func(t *testing.T) {
		executionCount := 0
		testFunc := func() {
			executionCount++
		}
		job := manager.skip(testFunc)
		require.NotNil(t, job)
		// 多次执行
		job.Run()
		job.Run()
		job.Run()
		assert.Equal(t, 3, executionCount)
	})
}

func TestSystemJobManager_delay(t *testing.T) {
	zapLogger := logger.NewLogger(zaptest.NewLogger(t).Sugar())
	emailClient := &initialize.EmailClient{}
	cronLogger := initialize.NewCronLogger(zapLogger, emailClient)
	cronClient := initialize.NewCronClient(cronLogger)
	manager := NewSystemJobManager(cronLogger, cronClient, zapLogger)

	t.Run("delay wrap func", func(t *testing.T) {
		executed := false
		testFunc := func() {
			executed = true
		}
		job := manager.delay(testFunc)
		require.NotNil(t, job)
		// 执行任务
		job.Run()
		assert.True(t, executed)
	})

	t.Run("delay func executed multiple times", func(t *testing.T) {
		executionCount := 0
		testFunc := func() {
			executionCount++
		}
		job := manager.delay(testFunc)
		require.NotNil(t, job)
		// 多次执行
		job.Run()
		job.Run()
		job.Run()
		assert.Equal(t, 3, executionCount)
	})
}

func TestSystemJobManager_RegisterJobsAndStart(t *testing.T) {
	t.Parallel()
	zapLogger := logger.NewLogger(zaptest.NewLogger(t).Sugar())
	emailClient := &initialize.EmailClient{}
	cronLogger := initialize.NewCronLogger(zapLogger, emailClient)
	cronClient := initialize.NewCronClient(cronLogger)

	t.Run("register for skip mode jobs", func(t *testing.T) {
		t.Parallel()
		mockJob := &MockSystemJob{
			name:           "test-skip-job",
			ifStillRunning: constant.Skip,
			interval:       "*/1 * * * * *", // 每秒执行一次
		}
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob)
		err := manager.RegisterJobsAndStart()
		require.NoError(t, err)
		// 验证任务已注册
		assert.Len(t, manager.jobsMap, 1)
		assert.Contains(t, manager.jobsMap, "test-skip-job")
		// 等待任务执行
		time.Sleep(1200 * time.Millisecond)
		assert.True(t, mockJob.handleCalled)
		// 清理
		manager.Stop()
	})

	t.Run("register for delay mode jobs", func(t *testing.T) {
		t.Parallel()
		mockJob := &MockSystemJob{
			name:           "test-delay-job",
			ifStillRunning: constant.Delay,
			interval:       "*/1 * * * * *", // 每秒执行一次
		}
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob)
		err := manager.RegisterJobsAndStart()
		require.NoError(t, err)
		// 验证任务已注册
		assert.Len(t, manager.jobsMap, 1)
		assert.Contains(t, manager.jobsMap, "test-delay-job")
		// 等待任务执行
		time.Sleep(2 * time.Second)
		assert.True(t, mockJob.handleCalled)
		// 清理
		manager.Stop()
	})

	t.Run("register for none mode jobs", func(t *testing.T) {
		t.Parallel()
		mockJob := &MockSystemJob{
			name:           "test-none-job",
			ifStillRunning: constant.None,
			interval:       "*/1 * * * * *", // 每秒执行一次
		}
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob)
		err := manager.RegisterJobsAndStart()
		require.NoError(t, err)
		// 验证任务已注册
		assert.Len(t, manager.jobsMap, 1)
		assert.Contains(t, manager.jobsMap, "test-none-job")
		// 等待任务执行
		time.Sleep(2 * time.Second)
		assert.True(t, mockJob.handleCalled)
		// 清理
		manager.Stop()
	})

	t.Run("register for multiple jobs", func(t *testing.T) {
		t.Parallel()
		mockJob1 := &MockSystemJob{
			name:           "test-job-1",
			ifStillRunning: constant.Skip,
			interval:       "*/1 * * * * *",
		}

		mockJob2 := &MockSystemJob{
			name:           "test-job-2",
			ifStillRunning: constant.Delay,
			interval:       "*/1 * * * * *",
		}

		mockJob3 := &MockSystemJob{
			name:           "test-job-3",
			ifStillRunning: constant.None,
			interval:       "*/1 * * * * *",
		}

		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob1, mockJob2, mockJob3)
		err := manager.RegisterJobsAndStart()
		require.NoError(t, err)
		// 验证所有任务都已注册
		assert.Len(t, manager.jobsMap, 3)
		assert.Contains(t, manager.jobsMap, "test-job-1")
		assert.Contains(t, manager.jobsMap, "test-job-2")
		assert.Contains(t, manager.jobsMap, "test-job-3")
		// 等待任务执行
		time.Sleep(2 * time.Second)
		assert.True(t, mockJob1.handleCalled)
		assert.True(t, mockJob2.handleCalled)
		assert.True(t, mockJob3.handleCalled)
		// 清理
		manager.Stop()
	})

	t.Run("invalid cron expression", func(t *testing.T) {
		t.Parallel()
		mockJob := &MockSystemJob{
			name:           "test-invalid-job",
			ifStillRunning: constant.Skip,
			interval:       "invalid cron expression", // 无效的cron表达式
		}
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob)
		err := manager.RegisterJobsAndStart()
		require.Error(t, err)
	})
}

func TestSystemJobManager_Stop(t *testing.T) {
	t.Parallel()
	zapLogger := logger.NewLogger(zaptest.NewLogger(t).Sugar())
	emailClient := &initialize.EmailClient{}
	cronLogger := initialize.NewCronLogger(zapLogger, emailClient)
	cronClient := initialize.NewCronClient(cronLogger)

	t.Run("stop with jobs SystemJobManager", func(t *testing.T) {
		t.Parallel()
		mockJob := &MockSystemJob{
			name:           "test-stop-job",
			ifStillRunning: constant.Skip,
			interval:       "*/1 * * * * *",
		}
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger, mockJob)
		// 注册并启动任务
		err := manager.RegisterJobsAndStart()
		require.NoError(t, err)
		// 验证任务已注册
		assert.Len(t, manager.jobsMap, 1)
		assert.Contains(t, manager.jobsMap, "test-stop-job")
		// 等待任务执行一次
		time.Sleep(1200 * time.Millisecond)
		assert.True(t, mockJob.handleCalled)
		// 记录执行次数
		executionCount := mockJob.handleCount
		// 停止
		manager.Stop()
		// 等待一段时间，验证任务不再执行
		time.Sleep(1 * time.Second)
		assert.Equal(t, executionCount, mockJob.handleCount)
	})

	t.Run("multiple calls to the Stop method", func(t *testing.T) {
		manager := NewSystemJobManager(cronLogger, cronClient, zapLogger)
		// 启动cron客户端
		cronClient.Start()
		// 多次停止（应该不会panic）
		manager.Stop()
		manager.Stop()
		manager.Stop()
	})
}
