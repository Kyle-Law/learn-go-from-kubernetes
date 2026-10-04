package basics

import (
	"reflect"
	"testing"
)

func readyPod() Pod {
	return Pod{
		Name: "web-1", Namespace: "default", Phase: PodRunning,
		Labels:     map[string]string{"app": "web", "tier": "frontend"},
		Containers: []Container{{Name: "app", Ready: true}, {Name: "sidecar", Ready: true}},
	}
}

func TestReady(t *testing.T) {
	notReadyContainer := readyPod()
	notReadyContainer.Containers[1].Ready = false

	pending := readyPod()
	pending.Phase = PodPending

	empty := readyPod()
	empty.Containers = nil

	tests := []struct {
		name string
		pod  Pod
		want bool
	}{
		{"all containers ready", readyPod(), true},
		{"one container not ready", notReadyContainer, false},
		{"pending phase", pending, false},
		{"no containers", empty, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pod.Ready(); got != tt.want {
				t.Errorf("Ready() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadyString(t *testing.T) {
	p := readyPod()
	p.Containers = append(p.Containers, Container{Name: "logger"})
	if got, want := p.ReadyString(), "2/3"; got != want {
		t.Errorf("ReadyString() = %q, want %q", got, want)
	}
	if got, want := (Pod{}).ReadyString(), "0/0"; got != want {
		t.Errorf("empty pod ReadyString() = %q, want %q", got, want)
	}
}

func TestSelectorMatches(t *testing.T) {
	labels := map[string]string{"app": "web", "tier": "frontend"}
	tests := []struct {
		name     string
		selector map[string]string
		labels   map[string]string
		want     bool
	}{
		{"exact match", map[string]string{"app": "web", "tier": "frontend"}, labels, true},
		{"subset match", map[string]string{"app": "web"}, labels, true},
		{"wrong value", map[string]string{"app": "db"}, labels, false},
		{"missing key", map[string]string{"env": "prod"}, labels, false},
		{"empty-string value vs missing key", map[string]string{"env": ""}, labels, false},
		{"nil selector matches all", nil, labels, true},
		{"empty selector matches nil labels", map[string]string{}, nil, true},
		{"selector vs nil labels", map[string]string{"app": "web"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SelectorMatches(tt.selector, tt.labels); got != tt.want {
				t.Errorf("SelectorMatches(%v, %v) = %v, want %v", tt.selector, tt.labels, got, tt.want)
			}
		})
	}
}

func TestFilterPods(t *testing.T) {
	pods := []Pod{
		{Name: "web-1", Namespace: "default", Labels: map[string]string{"app": "web"}},
		{Name: "db-1", Namespace: "default", Labels: map[string]string{"app": "db"}},
		{Name: "web-2", Namespace: "staging", Labels: map[string]string{"app": "web"}},
		{Name: "web-3", Namespace: "default", Labels: map[string]string{"app": "web"}},
	}
	names := func(ps []Pod) []string {
		var out []string
		for _, p := range ps {
			out = append(out, p.Name)
		}
		return out
	}
	tests := []struct {
		name      string
		namespace string
		selector  map[string]string
		want      []string
	}{
		{"namespace + selector", "default", map[string]string{"app": "web"}, []string{"web-1", "web-3"}},
		{"all namespaces", "", map[string]string{"app": "web"}, []string{"web-1", "web-2", "web-3"}},
		{"namespace only", "staging", nil, []string{"web-2"}},
		{"no matches", "prod", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := names(FilterPods(pods, tt.namespace, tt.selector)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterPods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountByPhase(t *testing.T) {
	pods := []Pod{{Phase: PodRunning}, {Phase: PodRunning}, {Phase: PodFailed}}
	want := map[Phase]int{PodRunning: 2, PodFailed: 1}
	if got := CountByPhase(pods); !reflect.DeepEqual(got, want) {
		t.Errorf("CountByPhase() = %v, want %v", got, want)
	}
}

func TestSetLabel(t *testing.T) {
	var p Pod // Labels is nil
	p.SetLabel("app", "web")
	p.SetLabel("app", "api")
	p.SetLabel("tier", "backend")
	want := map[string]string{"app": "api", "tier": "backend"}
	if !reflect.DeepEqual(p.Labels, want) {
		t.Errorf("Labels = %v, want %v", p.Labels, want)
	}
}
