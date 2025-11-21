package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, as *AuthService) {
	controller := &AuthController{AuthService: as}

	userGroup := r.Group("/course-tracker-api/user")

	{
		userGroup.POST("/login", controller.Login)
		userGroup.POST("/signup", controller.SignUp)
		userGroup.GET("/me", controller.Me)
	}

}
