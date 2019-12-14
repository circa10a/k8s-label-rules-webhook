package main

import "github.com/gin-gonic/gin"

func someHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}

	return gin.HandlerFunc(fn)
}

// Routes gin routes
func routes(router *gin.Engine) {
	root := router.Group("/")
	{
		root.GET("/", someHandler())

	}
}
