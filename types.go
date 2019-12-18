package main

// Structure of POST request from k8s admission webhook resource
type k8sRequest struct {
	APIVersion string
	Kind       string
	Request    struct {
		Object struct {
			Metadata struct {
				UID    string
				Labels map[string]interface{}
			}
		}
	}
}

// Structure of response to tell k8s if the resource is allowed or not
type webhookResponse struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Response   response `json:"response"`
}

type response struct {
	UID     string `json:"uid"`
	Allowed bool   `json:"allowed"`
	Status  status `json:"status"`
}

type status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Structs for other reponses within the API
type validRulesResponse struct {
	RulesValid bool        `json:"rulesValid"`
	Errors     []ruleError `json:"errors"`
}

type reloadRulesResponse struct {
	Reloaded   bool        `json:"reloaded"`
	YamlErr    string      `json:"yamlErr"`
	RuleErrs   []ruleError `json:"ruleErr"`
	NewRuleSet *[]rule     `json:"newRules"`
}
