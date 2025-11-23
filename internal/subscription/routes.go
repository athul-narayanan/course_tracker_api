package subscription

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, ss *SubscriptionService) {
	controller := &SubscriptionController{Service: ss}

	uni := r.Group("/course-tracker-api/subscription")
	{
		uni.POST("", controller.CreateSubscription)
	}
}
