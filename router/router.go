package router

import (
	"github.com/eclipse-aerios/llo-api/controllers"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// health := new(controllers.HealthController)
	health := &controllers.HealthController{}
	version := new(controllers.VersionController)

	router.GET("/health", health.Status)
	router.GET("/version", version.Version)
	// router.Use(middlewares.AuthMiddleware())

	v1 := router.Group("v1")
	{
		serviceComponentsGroup := v1.Group("service-components")
		{
			sc := new(controllers.ServiceComponentController)
			serviceComponentsGroup.GET("/", sc.List)
			// TODO add namespace as path parameter?
			// serviceComponentsGroup.GET("/:namespace", sc.List)
			serviceComponentsGroup.GET("/:scName", sc.GetByName)
			// serviceComponentsGroup.GET("/:namespace/:scName", sc.GetByName)
			serviceComponentsGroup.POST("/", sc.Deploy)
			serviceComponentsGroup.PUT("/", sc.Update)
			serviceComponentsGroup.PATCH("/:scName", sc.Patch)
			serviceComponentsGroup.DELETE("/:scName", sc.Delete)
		}
	}
	return router

}
