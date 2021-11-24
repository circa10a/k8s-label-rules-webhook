package main

// k8srequest is the object structure of a received POST request from k8s admission webhook resource
type k8sRequest struct {
	Request    request
	APIVersion string
	Kind       string
}

type request struct {
	Object object
}

type object struct {
	Metadata metadata
}

type metadata struct {
	Labels map[string]interface{}
	UID    string
}

// webhookResponse is the structure of a response to tell k8s if the resource is allowed or not
type webhookResponse struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Response   response `json:"response"`
}

type response struct {
	UID     string `json:"uid"`
	Status  status `json:"status"`
	Allowed bool   `json:"allowed"`
}

type status struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// validRulesResponse informs the user if there are any problems by accessing the /validate route
type validRulesResponse struct {
	Errors     []ruleError `json:"errors"`
	RulesValid bool        `json:"rulesValid"`
}

type reloadRulesResponse struct {
	NewRuleSet *[]rule     `json:"newRules"`
	YamlErr    string      `json:"yamlErr"`
	RuleErrs   []ruleError `json:"ruleErr"`
	Reloaded   bool        `json:"reloaded"`
}
