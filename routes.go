package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func labelValidationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		k8sData := k8sRequest{}
		if c.ShouldBind(&k8sData) == nil {
			log.Println(k8sData.Request.Object.Metadata.Labels)
		}

		c.String(200, "Success")
	}
	return gin.HandlerFunc(fn)
}

func reloadRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		err := R.load(*FilePath)
		if err != nil {
			c.PureJSON(400, gin.H{
				"err": err.Error(),
			})
		} else {
			c.PureJSON(200, gin.H{
				"newrules": &R,
			})
		}
	}

	return gin.HandlerFunc(fn)
}

func routes(router *gin.Engine) {
	root := router.Group("/")
	{
		root.POST("/", labelValidationHandler())
		root.POST("/reload", reloadRulesHandler())
	}
}
