package university

type University struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Field struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Specialization struct {
	ID      string `json:"id" gorm:"primaryKey"`
	FieldID string `json:"fieldId"`
	Name    string `json:"name"`
}

type Course struct {
	ID               uint `gorm:"primaryKey"`
	Name             string
	UniversityID     *int
	FieldID          *int
	SpecializationID *int
	Level            *string
	Duration         *string
	CourseLink       string
}

type CourseDTO struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Level          *string `json:"level"`
	Duration       *string `json:"duration"`
	CourseLink     string  `json:"courseLink"`
	University     string  `json:"university"`
	Field          string  `json:"field"`
	Specialization string  `json:"specialization"`
}

type Subscription struct {
	ID               uint   `gorm:"primaryKey"`
	Email            string `gorm:"column:user_email;index"`
	UniversityID     *int
	FieldID          *int
	SpecializationID *int
	Level            *string
	Duration         *string
}

func (Course) TableName() string {
	return "courses"
}

func (Subscription) TableName() string {
	return "subscriptions"
}
