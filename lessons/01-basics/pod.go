package basics

// Phase mirrors corev1.PodPhase.
type Phase string

const (
	PodPending   Phase = "Pending"
	PodRunning   Phase = "Running"
	PodSucceeded Phase = "Succeeded"
	PodFailed    Phase = "Failed"
)

type Container struct {
	Name  string
	Image string
	Ready bool
}

type Pod struct {
	Name       string
	Namespace  string
	Labels     map[string]string
	Phase      Phase
	Containers []Container
}

// Ready reports whether the pod is Running AND every container is ready.
// A pod with zero containers is not ready.
func (p Pod) Ready() bool {
	// TODO
	return false
}

// ReadyString returns the READY column of `kubectl get pods`: "<ready>/<total>",
// e.g. "1/2". Hint: fmt.Sprintf.
func (p Pod) ReadyString() string {
	// TODO
	return ""
}

// SelectorMatches reports whether labels contain every key=value pair in selector.
// An empty (or nil) selector matches everything — same as in Kubernetes.
func SelectorMatches(selector, labels map[string]string) bool {
	// TODO
	return false
}

// FilterPods returns the pods in namespace that match selector, preserving order.
// An empty namespace means "all namespaces" (like `kubectl get pods -A`).
func FilterPods(pods []Pod, namespace string, selector map[string]string) []Pod {
	// TODO
	return nil
}

// CountByPhase returns how many pods are in each phase.
// Phases with no pods should not appear in the map.
func CountByPhase(pods []Pod) map[Phase]int {
	// TODO
	return nil
}

// SetLabel sets a label on the pod, overwriting any existing value.
// It must work on a Pod whose Labels map is nil.
func (p *Pod) SetLabel(key, value string) {
	// TODO
}
