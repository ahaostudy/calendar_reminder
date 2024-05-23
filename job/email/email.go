package email

import (
	"bytes"
	"github.com/ahaostudy/calendar_reminder/conf"
	"github.com/ahaostudy/calendar_reminder/middleware/rabbitmq"
	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"html/template"
	"net/smtp"
)

const RMQEmailKey = "email"

var RMQEmail *rabbitmq.RabbitMQ

func RunEmailService() {
	RMQEmail = rabbitmq.NewWorkRabbitMQ(RMQEmailKey)
	go RMQEmail.Consume(sendEmailHandler)
}

// business functions that handle sending email requests
func sendEmailHandler(msg *amqp.Delivery) error {
	logrus.Info(msg.Body)

	return nil
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
	auth := smtp.PlainAuth("", conf.GetConf().Email.Email, conf.GetConf().Email.Auth, conf.GetConf().Email.Host)
	return e.Send(conf.GetConf().Email.Addr, auth)
}

type ReminderInfo struct {
	Title   string
	Content string
}

// HTML generate email HTML based on reminder info
func (ri *ReminderInfo) HTML() (string, error) {
	tmpl, err := template.ParseFiles("reminder.tmpl")
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, ri); err != nil {
		return "", err
	}
	return buf.String(), nil
}
