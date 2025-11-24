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

type CourseEvent struct {
	Name             string  `json:"name"`
	UniversityID     *int    `json:"universityId"`
	FieldID          *int    `json:"fieldId"`
	SpecializationID *int    `json:"specializationId"`
	Level            *string `json:"level"`
	Duration         *string `json:"duration"`
	Source           string  `json:"source"`
}

type Producer struct {
	writer *kafkago.Writer
	topic  string
}

func NewProducer() *Producer {
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

	log.Printf("Kafka producer using brokers: %v, topic: %s", brokers, topic)

	w := &kafkago.Writer{
		Addr:     kafkago.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}

	return &Producer{writer: w, topic: topic}
}

func (p *Producer) PublishCourseEvent(evt CourseEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	value, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	msg := kafkago.Message{
		Key:   []byte("course"),
		Value: value,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("Kafka publish error: %v", err)
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
