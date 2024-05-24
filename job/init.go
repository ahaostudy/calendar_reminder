package job

import (
	"github.com/ahaostudy/calendar_reminder/job/email"
	"github.com/ahaostudy/calendar_reminder/job/reminder"
)

func InitAsyncJobs() {
	email.RunEmailService()
	reminder.RunReminderService()
}

func DestroyAsyncJobs() {
	email.Destroy()
}
