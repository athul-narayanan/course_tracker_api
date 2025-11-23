package subscription

import (
	"course-tracker/config"

	"gorm.io/gorm"
)

type SubscriptionService struct {
	DB  *gorm.DB
	CFG *config.Config
}

func (s *SubscriptionService) CreateSubscription(sub Subscription) error {
	err := s.Save(sub)
	return err
}

func (r *SubscriptionService) Save(s Subscription) error {
	var exists int

	err := r.DB.Raw(`
		SELECT COUNT(*) FROM subscriptions 
		WHERE user_email = ? 
		  AND COALESCE(university_id,0) = COALESCE(?,0)
		  AND COALESCE(field_id,0) = COALESCE(?,0)
		  AND COALESCE(specialization_id,0) = COALESCE(?,0)
		  AND COALESCE(level,'') = COALESCE(?, '')
		  AND COALESCE(duration,'') = COALESCE(?, '')
	`,
		s.UserEmail,
		s.UniversityID,
		s.FieldID,
		s.SpecializationID,
		s.Level,
		s.Duration,
	).Scan(&exists).Error

	if err != nil {
		return err
	}

	if exists > 0 {
		return nil
	}

	return r.DB.Exec(`
		INSERT INTO subscriptions 
			(user_email, university_id, field_id, specialization_id, level, duration) 
		VALUES (?,?,?,?,?,?)
	`,
		s.UserEmail,
		s.UniversityID,
		s.FieldID,
		s.SpecializationID,
		s.Level,
		s.Duration,
	).Error
}
