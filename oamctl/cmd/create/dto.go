package create

type OamDeploy struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Components []Component `json:"components"`
}

type Component struct {
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	Properties ComponentProperties `json:"properties"`
	Traits     []Trait             `json:"traits"`
}

type ComponentProperties struct {
	Image string `json:"image"`
	Ports []Port `json:"ports"`
}

type Port struct {
	Port   int64 `json:"port"`
	Expose bool  `json:"expose"`
}

type Trait struct {
	Type       string          `json:"type"`
	Properties TraitProperties `json:"properties"`
}

type TraitProperties struct {
	Replicas int64 `json:"replicas"`
}
