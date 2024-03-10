package web

import "github.com/gin-gonic/gin"

func SetupProbeRoutes(r *gin.RouterGroup) {
	r.GET("/health/liveness", liveness)
	r.GET("/health/readiness", readiness)
	r.GET("/health/startup", startup)
}

func liveness(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func readiness(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func startup(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
