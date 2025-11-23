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
