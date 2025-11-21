package auth

type Auth struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string `gorm:"size:100;not null;column:firstname" json:"firstname"`
	LastName  string `gorm:"size:100;not null;column:lastname" json:"lastname"`
	Email     string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	Role      string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	ID          int    `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
	Role        string `json:"role"`
}

func (Auth) TableName() string {
	return "users"
}
