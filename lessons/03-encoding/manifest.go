package manifests

import (
	"errors"

	"sigs.k8s.io/yaml"
)

// Sentinel errors. Callers check them with errors.Is, even when wrapped.
var (
	ErrUnsupportedKind   = errors.New("unsupported kind")
	ErrContainerNotFound = errors.New("container not found")
	ErrMissingName       = errors.New("metadata.name is required")
	ErrNoContainers      = errors.New("at least one container is required")
	ErrMissingImage      = errors.New("image is required")
	ErrSelectorMismatch  = errors.New("selector does not match template labels")
)

// ParseDeployment parses a YAML or JSON manifest into a Deployment.
//
//   - If the input isn't valid YAML, return an error whose message starts with
//     "parsing deployment: " and wraps the original error (%w).
//   - If apiVersion/kind isn't exactly apps/v1 Deployment, return an error that
//     wraps ErrUnsupportedKind and says what it got, e.g.
//     "unsupported kind: v1/Service".
func ParseDeployment(data []byte) (*Deployment, error) {
	var d Deployment
	if err := yaml.Unmarshal(data, &d); err != nil {
		return nil, err // TODO: add context and wrap
	}
	// TODO: check apiVersion and kind
	return &d, nil
}

// ToJSON encodes the Deployment as compact JSON (no indentation).
// Hint: encoding/json.
func ToJSON(d *Deployment) ([]byte, error) {
	// TODO
	return nil, nil
}

// Ptr returns a pointer to a copy of v. It's generic: [T any] means it works
// for any type, so Ptr(int32(3)) is a *int32 and Ptr("x") is a *string.
// Kubernetes has the same helper: k8s.io/utils/ptr.To.
func Ptr[T any](v T) *T {
	// TODO
	return nil
}

// ReplicasOrDefault returns spec.replicas, or 1 if it isn't set
// (that's the default the API server applies).
func (d *Deployment) ReplicasOrDefault() int32 {
	// TODO
	return 0
}

// SetImage sets the image of the named container in the pod template — like
// `kubectl set image deployment/web app=nginx:1.28`. If there's no such
// container, return an error that wraps ErrContainerNotFound.
//
// Careful: `for _, c := range containers` gives you a COPY of each element.
func SetImage(d *Deployment, container, image string) error {
	// TODO
	return nil
}

// Validate checks the Deployment and returns ALL problems at once, combined
// with errors.Join (which returns nil if there are none):
//
//   - metadata.name is empty                -> ErrMissingName
//   - the pod template has no containers    -> ErrNoContainers
//   - a container has no image              -> ErrMissingImage, wrapped with the
//     container name: `container "app": image is required`
//   - selector.matchLabels is empty, or the template's labels don't satisfy
//     it (lesson 1's SelectorMatches!)       -> ErrSelectorMismatch
func Validate(d *Deployment) error {
	// TODO
	return nil
}
