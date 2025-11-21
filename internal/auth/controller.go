package auth

import (
	"course-tracker/config"
	"course-tracker/internal/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	AuthService *AuthService
}

func (ac *AuthController) SignUp(c *gin.Context) {
	var req struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password, err := util.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := Auth{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  password,
	}

	newuser, err := ac.AuthService.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"id":        newuser.ID,
			"firstname": newuser.FirstName,
			"lastname":  newuser.LastName,
			"email":     newuser.Email,
		},
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.AuthService.GetUser(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Oops! We couldn’t log you in. Please check your username and password and try again."})
		return
	}

	if err := util.VerifyPassword(req.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Oops! We couldn’t log you in. Please check your username and password and try again."})
		return
	}

	cfg := config.LoadConfig()

	tokenDuration := 30 * 24 * time.Hour
	tokenExp := time.Now().Add(tokenDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     tokenExp.Unix(),
	})

	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	httpOnly := true
	secure := false

	authCookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, authCookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": LoginResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
		},
	})
}

func (ac *AuthController) Me(c *gin.Context) {
	cfg := config.LoadConfig()

	accessToken, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
		return
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	user, err := ac.AuthService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": LoginResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
		},
	})
}
