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
