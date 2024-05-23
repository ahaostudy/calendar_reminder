package job

import "github.com/ahaostudy/calendar_reminder/job/email"

func InitAsyncJobs() {
	email.RunEmailService()
}

func DestroyAsyncJobs() {
	email.Destroy()
}
