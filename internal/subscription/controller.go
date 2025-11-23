package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	Service *SubscriptionService
}

func (c *SubscriptionController) CreateSubscription(ctx *gin.Context) {
	var req Subscription

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	if req.UserEmail == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User email required"})
		return
	}

	if err := c.Service.CreateSubscription(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Subscription saved"})
}
