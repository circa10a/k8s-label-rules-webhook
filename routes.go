package main

import (
	"log"

	"github.com/gin-gonic/gin"
	//"github.com/tidwall/gjson"
)

func someHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		k8sData := k8sRequest{}
		if c.ShouldBind(&k8sData) == nil {
			log.Println(k8sData.Request.Object.Metadata.Labels)
		}

		c.String(200, "Success")
	}

	return gin.HandlerFunc(fn)
}

func routes(router *gin.Engine) {
	root := router.Group("/")
	{
		root.POST("/", someHandler())

	}
}
