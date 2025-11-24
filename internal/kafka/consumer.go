package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type EmailService interface {
	SendCourseNotification(email string, evt CourseEvent, uniName string) error
}

type SubscriptionRepository interface {
	GetSubscribersForCourse(evt CourseEvent) ([]string, error)
	GetUniversityNameByID(id int) (string, error)
}

type Consumer struct {
	reader *kafkago.Reader
	repo   SubscriptionRepository
	email  EmailService
}

func NewConsumer(repo SubscriptionRepository, email EmailService) *Consumer {
	raw := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	var brokers []string
	if raw != "" {
		for _, p := range strings.Split(raw, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				brokers = append(brokers, p)
			}
		}
	}
	if len(brokers) == 0 {
		brokers = []string{"kafka:9092"}
	}

	topic := os.Getenv("KAFKA_COURSE_TOPIC")
	if topic == "" {
		topic = "course_updates"
	}

	groupID := os.Getenv("KAFKA_CONSUMER_GROUP")
	if groupID == "" {
		groupID = "course_notifier_group"
	}

	log.Printf("Kafka consumer using brokers: %v, topic: %s, group: %s", brokers, topic, groupID)

	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafkago.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait:     1 * time.Second,
	})

	return &Consumer{
		reader: r,
		repo:   repo,
		email:  email,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var evt CourseEvent
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			log.Printf("Failed to unmarshal CourseEvent: %v", err)
			continue
		}

		subscribers, err := c.repo.GetSubscribersForCourse(evt)
		if err != nil {
			log.Printf("Failed to load subscribers: %v", err)
			continue
		}

		universityName := ""
		if evt.UniversityID != nil {
			if name, err := c.repo.GetUniversityNameByID(*evt.UniversityID); err != nil {
				log.Printf("Failed to load university name for id=%d: %v", *evt.UniversityID, err)
			} else {
				universityName = name
			}
		}

		for _, email := range subscribers {
			if err := c.email.SendCourseNotification(email, evt, universityName); err != nil {
				log.Printf("Failed to send email to %s: %v", email, err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
