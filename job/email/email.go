package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/ahaostudy/calendar_reminder/model"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/ahaostudy/calendar_reminder/conf"
	"github.com/ahaostudy/calendar_reminder/middleware/rabbitmq"
)

const (
	RMQEmailKey   = "email"
	EmailAuthKey  = "EMAIL_AUTH"
	MaxRetryCount = 3
)

var RMQEmail *rabbitmq.RabbitMQ

func RunEmailService() {
	RMQEmail = rabbitmq.NewWorkRabbitMQ(RMQEmailKey)
	go RMQEmail.Consume(sendEmailHandler)
}

// business functions that handle sending email requests
func sendEmailHandler(msg *amqp.Delivery) {
	logrus.Info("send email request:", string(msg.Body))

	// unmarshal msg body
	m := new(Message)
	err := json.Unmarshal(msg.Body, m)
	if err != nil {
		logrus.Error(err)
		m.Send() // retry
		return
	}

	// send email
	task := m.Task
	t := time.Unix(task.Time, 0).Add(8 * time.Hour)
	ri := ReminderInfo{Title: task.Title, Time: t.Format("2006-01-02 15:04:05")}
	err = Send(task.Title, ri.HTML(), task.User.Email)
	if err != nil {
		logrus.Error(err)
		m.Send() // retry
		return
	}
}

// Send and retry automatically
func (msg *Message) Send() {
	if msg.RetryCount >= MaxRetryCount {
		return
	}
	msg.RetryCount++
	m, err := json.Marshal(msg)
	if err != nil {
		logrus.Error("json marshal error:", err.Error())
		return
	}
	err = RMQEmail.Publish(m)
	if err != nil {
		logrus.Error("publish to rabbitmq failed:", err.Error())
	}
}

type Message struct {
	RetryCount int         `json:"retry_count"`
	Task       *model.Task `json:"task"`
}

func Destroy() {
	RMQEmail.Destroy()
}

// Send email to a specified emails
func Send(subject, html string, toEmails ...string) error {
	e := email.NewEmail()
	e.From = conf.GetConf().Email.From
	e.To = toEmails
	e.Subject = subject
	e.HTML = []byte(html)
	auth := smtp.PlainAuth("", conf.GetConf().Email.Email, os.Getenv(EmailAuthKey), conf.GetConf().Email.Host)
	return e.Send(conf.GetConf().Email.Addr, auth)
}

type ReminderInfo struct {
	Title string
	Time  string
}

// HTML generate email HTML based on reminder info
func (ri *ReminderInfo) HTML() string {
	bakHTML := fmt.Sprintf(`<h1>日程提醒：%s</h1>`, ri.Title)

	tmpl, err := template.ParseFiles(filepath.Join("job", "email", "reminder.tmpl"))
	if err != nil {
		logrus.Error(err)
		return bakHTML
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, ri); err != nil {
		logrus.Error(err)
		return bakHTML
	}
	return buf.String()
}
