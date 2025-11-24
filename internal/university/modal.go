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

type CourseDetail struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	University     string `json:"university"`
	Field          string `json:"field"`
	Specialization string `json:"specialization"`
	Level          string `json:"level"`
	Duration       string `json:"duration"`
	CourseLink     string `json:"courseLink"`
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

type Subscription struct {
	ID               uint   `gorm:"primaryKey"`
	Email            string `gorm:"index"`
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
