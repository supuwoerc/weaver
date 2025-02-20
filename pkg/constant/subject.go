package constant

type Subject string

const (
	Signup      Subject = "Signup"
	CronRecover Subject = "Cron Recover"
	Recover     Subject = "Recover"
	UnlockFail  Subject = "Redis Unlock Fail"
)
