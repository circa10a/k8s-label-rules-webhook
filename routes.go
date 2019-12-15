package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func sendResponse(c *gin.Context, code int, uid string, allowed bool, message string) {
	r := &webhookResponse{
		APIVersion: "admission.k8s.io/v1beta1",
		Kind:       "AdmissionsReview",
		Response: response{
			UID:     uid,
			Allowed: allowed,
			Status: status{
				Code:    code,
				Message: message,
			},
		},
	}
	c.JSON(code, &r)
}

func labelValidationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		k8sData := &k8sRequest{}
		// If no error Binding
		if c.ShouldBindJSON(&k8sData) == nil {
			labels := k8sData.Request.Object.Metadata.Labels
			uid := k8sData.Request.Object.Metadata.UID
			// Loop current ruleset
			err := ensureLabelsContainRules(labels)
			if err != nil {
				sendResponse(c, http.StatusBadRequest, uid, false, err.Error())
			}
			//c.String(200, "Success")
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "Improperly formatted request sent",
			})
		}
	}
	return gin.HandlerFunc(fn)
}

func reloadRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		err := R.load(*FilePath)
		if err != nil {
			c.PureJSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
		} else {
			c.PureJSON(http.StatusOK, gin.H{
				"newrules": &R,
			})
		}
	}
	return gin.HandlerFunc(fn)
}

func getRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.PureJSON(http.StatusOK, &R)
	}
	return gin.HandlerFunc(fn)
}

func routes(router *gin.Engine) {
	root := router.Group("/")
	{
		root.POST("/", labelValidationHandler())
		root.POST("/reload", reloadRulesHandler())
		root.GET("/rules", getRulesHandler())
	}
}
