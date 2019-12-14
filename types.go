package main

// Rules array from yaml
type rules struct {
	Rules []rule `yaml:"rules"`
}

// Individual rule within rules array
type rule struct {
	Name  string `yaml:"name"`
	Key   string `yaml:"key"`
	Value value  `yaml:"value"`
}

// Value struct within each rule
type value struct {
	DataType string `yaml:"type"`
	Regex    string `yaml:"regex"`
}

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
