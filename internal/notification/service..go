package notification

import (
	"course-tracker/config"
	"fmt"

	"gorm.io/gorm"
)

type NotificationService struct {
	DB     *gorm.DB
	Config *config.Config
}

func (s *NotificationService) CreateNotification(email string, message string) error {
	notification := Notification{
		Email:   email,
		Message: message,
		Status:  "new",
	}

	if err := s.DB.Create(&notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil

}

func (s *NotificationService) MarkAsRead(notificationID uint) error {
	return s.DB.Model(&Notification{}).
		Where("id = ?", notificationID).
		Update("status", "read").Error
}

func (s *NotificationService) GetNotificationsForUser(email string) ([]Notification, error) {
	var notifications []Notification
	if err := s.DB.Where("email = ?", email).Order("created_at desc").Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
