package subscription

type Subscription struct {
	ID               int     `json:"id,omitempty"`
	UserEmail        string  `json:"userEmail" binding:"required,email"`
	UniversityID     *int    `json:"universityId"`
	FieldID          *int    `json:"fieldId"`
	SpecializationID *int    `json:"specializationId"`
	Level            *string `json:"level"`
	Duration         *string `json:"duration"`
	CreatedAt        string  `json:"createdAt,omitempty"`
}
