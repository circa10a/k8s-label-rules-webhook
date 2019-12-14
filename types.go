package main

// Structure of POST request from k8s admission webhook resource
type k8sRequest struct {
	Request struct {
		Object struct {
			Metadata struct {
				Labels map[string]interface{}
			}
		}
	}
}
