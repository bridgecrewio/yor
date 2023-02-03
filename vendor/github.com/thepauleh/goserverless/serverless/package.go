package serverless

// Package definition
type Package struct {
	Artifact               string   `json:"artifact,omitempty"`
	ExcludeDevDependencies bool     `json:"excludeDevDependencies,omitempty"`
	Exclude                []string `json:"exclude,omitempty"`
	Individually           bool     `json:"individually,omitempty"`
	Include                []string `json:"include,omitempty"`
}
