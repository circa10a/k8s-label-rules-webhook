package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Generate response to return to k8s
func sendResponse(c *gin.Context, k8sRequest *k8sRequest, uid string, allowed bool, code int, message string) {
	r := &webhookResponse{
		APIVersion: k8sRequest.ApiVersion,
		Kind:       k8sRequest.Kind,
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
		// If no binding err
		if c.BindJSON(&k8sData) == nil {
			labels := k8sData.Request.Object.Metadata.Labels
			uid := k8sData.Request.Object.Metadata.UID
			// First check ruleset is valid
			rulesErr := validateAllRulesRegex(R)
			if len(rulesErr) > 0 {
				sendResponse(c, k8sData, uid, false, http.StatusInternalServerError, strings.Join(rulesErr[:], ","))
				return
			}
			// Ensure labels provided contain keys identified in the ruleset
			containLabelErr := ensureLabelsContainRules(labels)
			// Reject request if not
			if containLabelErr != nil {
				sendResponse(c, k8sData, uid, false, http.StatusBadRequest, containLabelErr.Error())
				return
			}
			// Ensure labels provided match regex of keys identified in the ruleset
			matchLabelErr := ensureLabelsMatchRules(labels)
			// Reject request if not
			if matchLabelErr != nil {
				sendResponse(c, k8sData, uid, false, http.StatusBadRequest, matchLabelErr.Error())
				return
			}
			// All constraints passed
			sendResponse(c, k8sData, uid, true, http.StatusOK, "Labels conform to ruleset")
			return
		}
		// In the event, nothing to bind to
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Improperly formatted request sent",
		})
		return
	}
	return gin.HandlerFunc(fn)
}

// /reload context
// Hot reload file back into memory via pointer
func reloadRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Load file back into mem
		err := R.load(*FilePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"reloaded": false,
				"err":      err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"reloaded":   true,
				"newRuleSet": &R,
			})
		}
	}
	return gin.HandlerFunc(fn)
}

// /rules context
// Send current ruleset
func getRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, &R)
	}
	return gin.HandlerFunc(fn)
}

// /validate context
// Send whether ruleset regex is valid or not
func validateRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		err := validateAllRulesRegex(R)
		if len(err) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"rulesValid": false,
				"err":        err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
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
