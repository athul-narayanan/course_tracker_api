package university

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, us *UniversityService) {
	controller := &UniversityController{UniversityService: us}

	uni := r.Group("/course-tracker-api")
	{
		uni.GET("/universities", controller.GetUniversities)
		uni.GET("/fields", controller.GetFields)
		uni.GET("/specializations", controller.GetSpecializations)
		uni.GET("/universities/search", controller.SearchUniversities)
		uni.POST("/universities/add", controller.AddCourse)
	}
}
