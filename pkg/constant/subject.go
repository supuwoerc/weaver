package constant

type Subject string

const (
	Signup      Subject = "Signup"
	CronRecover Subject = "Cron Recover"
	Recover     Subject = "Recover"
	UnlockFail  Subject = "Redis Unlock Fail"
	ExtendErr   Subject = "Redis Lock Extend Err"
	ExtendFail  Subject = "Redis Lock Extend Fail"
)
