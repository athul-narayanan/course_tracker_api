package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService *NotificationService
}

func (c *NotificationController) GetNotifications(ctx *gin.Context) {
	email := ctx.Query("email")

	notifications, err := c.NotificationService.GetNotificationsForUser(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notifications"})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	if err := c.NotificationService.MarkAsRead(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark notification as read"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}

type createNotificationRequest struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (c *NotificationController) CreateNotification(ctx *gin.Context) {
	var req createNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := c.NotificationService.CreateNotification(req.Email, req.Message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "notification created"})

}
