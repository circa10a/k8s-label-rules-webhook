package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Generate response to return to k8s
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

// / context
// K8s will call this handler for rule validation
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

// /reload context
func reloadRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Load file back into mem
		err := R.load(*FilePath)
		if err != nil {
			c.PureJSON(http.StatusBadRequest, gin.H{
				"reloaded": false,
				"err":      err.Error(),
			})
		} else {
			c.PureJSON(http.StatusOK, gin.H{
				"reloaded":   true,
				"newRuleSet": &R,
			})
		}
	}
	return gin.HandlerFunc(fn)
}

// /rules context
func getRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.PureJSON(http.StatusOK, &R)
	}
	return gin.HandlerFunc(fn)
}

// /validate context
func validateRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		err := validateAllRulesRegex(R)
		if len(err) > 0 {
			c.PureJSON(http.StatusOK, gin.H{
				"rulesValid": false,
				"err":        err,
			})
		} else {
			c.PureJSON(http.StatusOK, gin.H{
				"rulesValid": true,
			})
		}
	}
	return gin.HandlerFunc(fn)
}

// Init all context paths
func routes(router *gin.Engine) {
	root := router.Group("/")
	{
		// K8s calls for validation
		root.POST("/", labelValidationHandler())
		// Hot reload rules.yaml file back into memory
		// Allows updating of rules without restarting of application
		root.POST("/reload", reloadRulesHandler())
		// View all current rules loaded into memory
		root.GET("/rules", getRulesHandler())
		// Validate regex pattern of rules
		root.GET("/validate", validateRulesHandler())
	}
}
