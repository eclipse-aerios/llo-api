package controllers

import (
	"net/http"

	"github.com/eclipse-aerios/llo-api/config"
	"github.com/eclipse-aerios/llo-api/models"

	"github.com/gin-gonic/gin"
)

var buildTime string
var commitHash string

type VersionController struct{}

func (v VersionController) Version(c *gin.Context) {
	apiVersion := models.ApiVersion{
		Version:      config.API_VERSION,
		BuildTime:    buildTime,
		CommitHash:   commitHash,
		ServiceName:  config.SERVICE_NAME,
		SupportedCRs: config.GetSupportedCRs(),
	}
	c.JSON(http.StatusOK, apiVersion)
}
