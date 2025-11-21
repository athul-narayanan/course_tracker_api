package auth

import (
	"course-tracker/config"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type AuthService struct {
	DB  *gorm.DB
	CFG *config.Config
}

func (s *AuthService) CreateUser(user Auth) (*Auth, error) {
	if user.Role == "" {
		user.Role = "User"
	}

	if err := s.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("An account with this email or phone number already exists. Please log in or use different details.")
		}
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) GetUser(email string) (*Auth, error) {

	var user Auth
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *AuthService) GetUserByID(id int) (*Auth, error) {
	var user Auth
	result := s.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
