package exporter

const (
	DefaultClusterName = "default"

	ComponentLabelKey = "app.kubernetes.io/component"
	ClusterLabelKey   = "app.kubernetes.io/cluster"

	// FQNamespace represents the prometheus FQName
	FQNamespace  = "nebula"
	NonNamespace = "none_namespace"
)

type (
	StaticConfig struct {
		Version  string    `yaml:"version"`
		Clusters []Cluster `yaml:"clusters"`
		ExtraLabels []Label `yaml:"extraLabels"`
	}

	Label struct {
		Name string `yaml:"name"`
		Value string `yaml:"value"`
	}

	Cluster struct {
		Name      string     `yaml:"name"`
		Instances []Instance `yaml:"instances"`
	}

	Instance struct {
		Name          string `yaml:"name"`
		EndpointIP    string `yaml:"endpointIP"`
		EndpointPort  int32  `yaml:"endpointPort"`
		ComponentType string `yaml:"componentType"`
	}
)
