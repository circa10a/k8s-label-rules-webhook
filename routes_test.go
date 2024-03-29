package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect"
	log "github.com/sirupsen/logrus"
)

func init() {
	// set defaults to global pointers for config
	rulesFile := "rules.yaml"
	filePath = &rulesFile
	// Init map to store compiled regexs
	r.compiledRegexs = make(map[string]*regexp.Regexp)
	// Load initial rules into memory
	err := r.load(*filePath)
	if err != nil {
		log.Error(err)
	}
	// load handlers into gin engine
	routes(g)
}

func TestRulesEndpoint(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)
	// is it working?
	response := e.GET("/rules").
		Expect().
		Status(http.StatusOK).JSON().Array()
	// Ensure correct count of rules
	response.Length().Equal(len(r.Rules))
	// Ensure no data is malformed from yaml to response
	for i := range response.Iter() {
		response.Element(i).Object().ValueEqual("name", r.Rules[i].Name)
		response.Element(i).Object().ValueEqual("key", r.Rules[i].Key)
		response.Element(i).Object().Value("value").Object().ValueEqual("regex", r.Rules[i].Value.Regex)
	}
}

func TestReloadEndpoint(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)
	// is it working?
	response := e.POST("/reload").
		Expect().
		Status(http.StatusOK).JSON()

	response.Object().ValueEqual("reloaded", true)
	response.Object().ValueEqual("yamlErr", "")
	response.Object().ValueEqual("ruleErr", nil)
	newRulesResponse := response.Object().Value("newRules").Array()
	// Ensure correct count of rules
	newRulesResponse.Length().Equal(len(r.Rules))
	// Ensure no data is malformed from yaml to response
	for i := range newRulesResponse.Iter() {
		newRulesResponse.Element(i).Object().ValueEqual("name", r.Rules[i].Name)
		newRulesResponse.Element(i).Object().ValueEqual("key", r.Rules[i].Key)
		newRulesResponse.Element(i).Object().Value("value").Object().ValueEqual("regex", r.Rules[i].Value.Regex)
	}
}

func TestValidateEndpoint(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)
	// is it working?
	response := e.GET("/validate").
		Expect().
		Status(http.StatusOK).JSON()

	response.Object().ValueEqual("rulesValid", true)
	response.Object().ValueEqual("errors", nil)
}

func TestRootEndpointNoMatchLabels(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	request := &k8sRequest{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Request: request{
			Object: object{
				Metadata: metadata{
					UID: "123",
					Labels: map[string]interface{}{
						"test": "value",
					},
				},
			},
		},
	}
	// is it working?
	response := e.POST("/").WithJSON(request).
		Expect().
		Status(http.StatusOK).JSON()

	// Validate response to k8s
	response.Object().ValueEqual("apiVersion", "admission.k8s.io/v1")
	response.Object().ValueEqual("kind", "AdmissionReview")
	response.Object().Value("response").Object().ValueEqual("allowed", false)
	response.Object().Value("response").Object().ValueEqual("uid", "123")
	response.Object().Value("response").Object().Value("status").Object().ValueEqual("code", 403)
	//nolint
	response.Object().Value("response").Object().Value("status").Object().ValueEqual("message", "phone-number not in labels")
}

// nolint
func TestRootEndpointLabelsInvalidRegex(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	request := &k8sRequest{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Request: request{
			Object: object{
				Metadata: metadata{
					UID: "123",
					Labels: map[string]interface{}{
						"phone-number": "value",
						"number":       "0",
					},
				},
			},
		},
	}
	// is it working?
	response := e.POST("/").WithJSON(request).
		Expect().
		Status(http.StatusOK).JSON()

	// Validate response to k8s
	response.Object().ValueEqual("apiVersion", "admission.k8s.io/v1")
	response.Object().ValueEqual("kind", "AdmissionReview")
	response.Object().Value("response").Object().ValueEqual("allowed", false)
	response.Object().Value("response").Object().ValueEqual("uid", "123")
	response.Object().Value("response").Object().Value("status").Object().ValueEqual("code", 403)
	response.Object().Value("response").Object().Value("status").Object().ValueEqual("message", "Value for label 'phone-number' does not match expression '[0-9]{3}-[0-9]{3}-[0-9]{4}'")
}

// nolint
func TestRootEndpointValidLabels(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	request := &k8sRequest{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Request: request{
			Object: object{
				Metadata: metadata{
					UID: "123",
					Labels: map[string]interface{}{
						"phone-number": "555-555-5555",
						"number":       0,
					},
				},
			},
		},
	}
	// is it working?
	response := e.POST("/").WithJSON(request).
		Expect().
		Status(http.StatusOK).JSON()

	// Validate response to k8s
	response.Object().ValueEqual("apiVersion", "admission.k8s.io/v1")
	response.Object().ValueEqual("kind", "AdmissionReview")
	response.Object().Value("response").Object().ValueEqual("allowed", true)
	response.Object().Value("response").Object().ValueEqual("uid", "123")
	response.Object().Value("response").Object().Value("status").Object().ValueEqual("code", 200)
	response.Object().Value("response").Object().Value("status").Object().ValueEqual("message", "Labels conform to ruleset")
}

func TestRootEndpointNoPayload(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// is it working?
	response := e.POST("/").WithJSON("").
		Expect().
		Status(http.StatusBadRequest).JSON()
	// Validate err response to empty body POST
	response.Object().ValueEqual("err", "Improperly formatted request sent")
}

func TestUndefinedRouteRedirect(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// is it working?
	response := e.GET("/notthere").
		Expect().
		Status(http.StatusOK)
	// Ensure proper redirect to swagger docs
	response.Body().Contains("swagger-ui")
}

func TestMetricsEndpoint(t *testing.T) {
	t.Parallel()
	// run server using httptest
	server := httptest.NewServer(g)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// is it working?
	response := e.GET("/metrics").
		Expect().
		Status(http.StatusOK).ContentType("text/plain")
	// Ensure proper redirect to swagger docs
	response.Body().Contains("gin_request_duration_seconds_sum")
}
