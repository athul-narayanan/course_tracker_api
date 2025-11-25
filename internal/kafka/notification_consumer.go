package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Notification struct {
	UserID    int
	Message   string
	Status    string
	CreatedAt time.Time
}

type NotificationService interface {
	CreateNotification(email string, message string) error
}

type NotificationConsumer struct {
	reader   *kafkago.Reader
	repo     SubscriptionRepository
	notifier NotificationService
}

func NewNotificationConsumer(repo SubscriptionRepository, notifier NotificationService) *NotificationConsumer {
	raw := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	mode := os.Getenv("KAFKA_CONSUMER_STRATEGY")

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

	groupID := os.Getenv("KAFKA_NOTIFICATION_GROUP")
	if groupID == "" {
		groupID = "course_notification_group"
	}

	log.Printf("Kafka notification consumer using brokers: %v, topic: %s, group: %s, strategy=%s",
		brokers, topic, groupID, mode)

	cfg := kafkago.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafkago.FirstOffset,
	}

	if mode == "fixed" {
		cfg.MinBytes = 10e3
		cfg.MaxBytes = 10e6
		cfg.MaxWait = 1 * time.Second
	} else {
		cfg.MinBytes = 1
		cfg.MaxBytes = 100
		cfg.MaxWait = 100 * time.Millisecond
	}

	r := kafkago.NewReader(cfg)

	return &NotificationConsumer{
		reader:   r,
		repo:     repo,
		notifier: notifier,
	}
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
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

		universityName := ""
		if evt.UniversityID != nil {
			if name, err := c.repo.GetUniversityNameByID(*evt.UniversityID); err != nil {
				log.Printf("Failed to load university name for id=%d: %v", *evt.UniversityID, err)
			} else {
				universityName = name
			}
		}

		userIDs, err := c.repo.GetSubscribersForCourse(evt)
		if err != nil {
			log.Printf("Failed to get subscribers for course: %v", err)
			continue
		}

		for _, userID := range userIDs {
			body := fmt.Sprintf(
				"Hi,\n\nA new course that matches your preferences has been added at %s:\n\nName: %s\nLevel: %s\nDuration: %s\n\nRegards,\nCourse Tracker\n\nVisit the link to see details: %s",
				universityName,
				evt.Name,
				*evt.Level,
				*evt.Duration,
				evt.CourseLink,
			)

			if err := c.notifier.CreateNotification(userID, body); err != nil {
				log.Printf("Failed to create notification for user %d: %v", userID, err)
			}
		}
	}

}

func (c *NotificationConsumer) Close() error {
	return c.reader.Close()
}
