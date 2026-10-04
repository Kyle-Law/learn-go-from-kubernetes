package manifests

// These types mirror a cut-down apps/v1 Deployment.
//
// TODO: add `json:"..."` struct tags to every field so these types read AND
// write real Kubernetes YAML/JSON (camelCase keys, empty fields omitted where
// noted). One tag is done for you. sigs.k8s.io/yaml converts YAML to JSON
// first, so json tags are all you need — exactly like the real k8s.io/api types.

type TypeMeta struct {
	APIVersion string `json:"apiVersion"`
	Kind       string
}

type ObjectMeta struct {
	Name        string            // omit when empty (pod templates have no name)
	Namespace   string            // omit when empty
	Labels      map[string]string // omit when empty
	Annotations map[string]string // omit when empty
}

type Deployment struct {
	// Embedded with no tag: its fields (apiVersion, kind) appear at the top
	// level of the JSON. Real k8s types write `json:",inline"` here, which
	// means the same thing to encoding/json.
	TypeMeta
	Metadata ObjectMeta
	Spec     DeploymentSpec
	Status   *DeploymentStatus // omit when nil
}

type DeploymentSpec struct {
	// A pointer, so "not set" (nil) is different from "set to 0".
	Replicas *int32 // omit when nil
	Selector LabelSelector
	Template PodTemplateSpec
}

type LabelSelector struct {
	MatchLabels map[string]string // omit when empty
}

type PodTemplateSpec struct {
	Metadata ObjectMeta
	Spec     PodSpec
}

type PodSpec struct {
	Containers []Container
}

type Container struct {
	Name  string
	Image string
	Ports []ContainerPort // omit when empty
	Env   []EnvVar        // omit when empty
}

type ContainerPort struct {
	ContainerPort int32
	Protocol      string // omit when empty
}

type EnvVar struct {
	Name  string
	Value string // omit when empty
}

type DeploymentStatus struct {
	Replicas      int32 // omit when zero
	ReadyReplicas int32 // omit when zero
}
