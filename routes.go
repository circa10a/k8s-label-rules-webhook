package main

import (
	"net/http"

	_ "github.com/circa10a/k8s-label-rules-webhook/api"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// K8s webhook context godoc
// @Summary Respond to k8s with approval or rejection of labels compared against the ruleset
// @Description Respond to k8s with approval or rejection of labels compared against the ruleset
// @Produce json
// @Success 200 {object} webhookResponse
// @Router // [post]
func labelValidationHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		k8sData := &k8sRequest{}
		// If no binding err
		if c.ShouldBindJSON(&k8sData) == nil {
			labels := k8sData.Request.Object.Metadata.Labels
			uid := k8sData.Request.Object.Metadata.UID
			// Ensure labels contain keys identified in the ruleset
			// Ensure labels are present and user provided regex match values identified in the ruleset
			matchLabelErr := r.ensureLabelsMatchRules(labels)
			// Reject request if not
			if matchLabelErr != nil {
				c.JSON(http.StatusOK, &webhookResponse{
					APIVersion: k8sData.APIVersion,
					Kind:       k8sData.Kind,
					Response: response{
						UID:     uid,
						Allowed: false,
						Status: status{
							Code:    http.StatusForbidden,
							Message: matchLabelErr.Error(),
						},
					},
				})
				return
			}
			c.JSON(http.StatusOK, &webhookResponse{
				APIVersion: k8sData.APIVersion,
				Kind:       k8sData.Kind,
				Response: response{
					UID:     uid,
					Allowed: true,
					Status: status{
						Code:    http.StatusOK,
						Message: "Labels conform to ruleset",
					},
				},
			})
			return
		}
		// In the event, nothing to bind to
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Improperly formatted request sent",
		})
	}
	return gin.HandlerFunc(fn)
}

// Reload godoc
// @Summary Reload ruleset file
// @Description Hot reload ruleset file into memory without downtime
// @Produce json
// @Success 200 {array} reloadRulesResponse
// @Router /reload [post]
func reloadRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Load file back into mem
		yamlErr := r.load(*filePath)
		// First check ruleset is valid
		ruleErrs := r.validateAllRulesRegex()
		reloadRulesResponse := &reloadRulesResponse{
			Reloaded:   true,
			YamlErr:    errToStr(yamlErr),
			RuleErrs:   ruleErrs,
			NewRuleSet: &r.Rules,
		}
		c.JSON(http.StatusOK, &reloadRulesResponse)
	}
	return gin.HandlerFunc(fn)
}

// Rules godoc
// @Summary Display loaded rules
// @Description Show current ruleset being used to validate labels against
// @Produce json
// @Success 200 {array} rule
// @Router /rules [get]
func getRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, &r.Rules)
	}
	return gin.HandlerFunc(fn)
}

// Validate godoc
// @Summary Validate regex of all loaded rules
// @Description Validate regex of every rule in ruleset, respond with rule errors
// @Produce json
// @Success 200 {object} validRulesResponse
// @Router /validate [get]
func validateRulesHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var rulesValid bool
		err := r.validateAllRulesRegex()
		if err != nil {
			rulesValid = false
		} else {
			rulesValid = true
		}
		c.JSON(http.StatusOK, &validRulesResponse{RulesValid: rulesValid, Errors: err})
	}
	return gin.HandlerFunc(fn)
}

func noRouteHandler() gin.HandlerFunc {
	// Replace doc.json to index.html to ensure user is brought to swagger site
	fn := func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
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
		// Swagger API Docs
		root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	}
	// For every route not defined, forward to swagger docs
	router.NoRoute(noRouteHandler())
}
