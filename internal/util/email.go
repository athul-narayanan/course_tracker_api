package util

import (
	"course-tracker/internal/kafka"
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (s *EmailService) SendCourseNotification(to string, evt kafka.CourseEvent, uniName string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	level := ""
	if evt.Level != nil {
		level = *evt.Level
	}
	duration := ""
	if evt.Duration != nil {
		duration = *evt.Duration
	}

	subject := "New course matching your preferences"
	body := fmt.Sprintf(
		"Hi,\n\nA new course that matches your preferences has been added at %s:\n\nName: %s\nLevel: %s\nDuration: %s\n\nRegards,\nCourse Tracker\n",
		uniName,
		evt.Name,
		level,
		duration,
	)

	msg := "From: " + s.from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
