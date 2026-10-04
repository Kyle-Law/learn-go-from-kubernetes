package basics

import "fmt"

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
	if p.Phase != PodRunning || len(p.Containers) == 0 {
		return false
	}
	for _, c := range p.Containers {
		if !c.Ready {
			return false
		}
	}
	return true
}

// ReadyString returns the READY column of `kubectl get pods`: "<ready>/<total>",
// e.g. "1/2". Hint: fmt.Sprintf.
func (p Pod) ReadyString() string {
	// TODO
	// PSEUDO CODE
	// p is a Pod struct, which container Publicly available attributes like Containers, ... etc
	// ready / len(p.Containers)
	// ready :=0
	// Loop through p.Containers, if c.Ready -> ready++
	ready := 0
	for _, c := range p.Containers {
		if c.Ready {
			ready++
		}
	}
	return fmt.Sprintf("%d/%d", ready, len(p.Containers))
}

// SelectorMatches reports whether labels contain every key=value pair in selector.
// An empty (or nil) selector matches everything — same as in Kubernetes.
func SelectorMatches(selector, labels map[string]string) bool {
	// Above is a shorthand for selector map[string]string, labels map[string]string
	// TODO
	// PSEUDOCODE
	// Loop over K,V of Selector
	// how to loop?
	// Ans: for k,v := range(selector) {}
	// if labels[k] not exists, or labels[k] != V -> return false
	for k, v := range selector {
		got, exist := labels[k]
		if !exist || got != v {
			return false
		}
	}
	return true
}

// FilterPods returns the pods in namespace that match selector, preserving order.
// An empty namespace means "all namespaces" (like `kubectl get pods -A`).
// Kyle: so basically here it means to implement sth like:
// k get po -n <namespace> -l a=b, c=d
func FilterPods(pods []Pod, namespace string, selector map[string]string) []Pod {
	// TODO
	// PSEUDOCODE
	var filteredPods []Pod
	// Loop over each pod, continue if namespace == "" OR pod.Namespace != namespace
	// Loop over selector - k,v
	// if [pod.Labels] exists and == v
	// filteredPods.add(pod)
	for _, p := range pods {
		if namespace == "" || p.Namespace == namespace {
			if SelectorMatches(selector, p.Labels) {
				filteredPods = append(filteredPods, p)
			}
		}
	}
	return filteredPods
}

// CountByPhase returns how many pods are in each phase.
// Phases with no pods should not appear in the map.
func CountByPhase(pods []Pod) map[Phase]int {
	// TODO
	phaseCount := make(map[Phase]int)

	for _, p := range pods {
		phaseCount[p.Phase]++
	}
	return phaseCount
}

// SetLabel sets a label on the pod, overwriting any existing value.
// It must work on a Pod whose Labels map is nil.
func (p *Pod) SetLabel(key, value string) {
	if p.Labels == nil {
		p.Labels = make(map[string]string)
	}
	p.Labels[key] = value
}
