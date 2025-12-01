package controllers

import (
	"net/http"

	"github.com/eclipse-aerios/llo-api/config"
	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	if config.Status == config.HEALTHY_STATUS {
		c.JSON(http.StatusOK, gin.H{"message": "The aeriOS LLO API is running"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "The aeriOS LLO API is UNHEALTHY"})
	}
}
