package notification

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, ns *NotificationService) {
	controller := &NotificationController{NotificationService: ns}

	uni := r.Group("/course-tracker-api/notifications")
	{
		uni.POST("", controller.CreateNotification)
		uni.GET("", controller.GetNotifications)
		uni.GET("/read", controller.MarkAsRead)
	}
}
