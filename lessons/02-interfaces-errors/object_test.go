package objects

import "testing"

// Compile-time interface checks. If *Pod stopped satisfying Object, the
// package would not build. You'll see this idiom all over Kubernetes.
var (
	_ Object = &Pod{}
	_ Object = &ConfigMap{}
	_ error  = &StatusError{}
)

// Service is defined only in this test file, yet it is an Object because it
// has the right methods. Interfaces are satisfied implicitly.
type Service struct {
	ObjectMeta
	Port int
}

func (s *Service) GetKind() string { return "Service" }

func TestPromotedMethods(t *testing.T) {
	var obj Object = &Pod{ObjectMeta: ObjectMeta{Name: "web", Namespace: "default"}}
	if got := obj.GetName(); got != "web" {
		t.Errorf("GetName() = %q, want %q", got, "web")
	}
	if got := obj.GetNamespace(); got != "default" {
		t.Errorf("GetNamespace() = %q, want %q", got, "default")
	}
	if got := obj.GetKind(); got != "Pod" {
		t.Errorf("GetKind() = %q, want %q", got, "Pod")
	}
	if got := (&ConfigMap{}).GetKind(); got != "ConfigMap" {
		t.Errorf("ConfigMap GetKind() = %q, want %q", got, "ConfigMap")
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		name string
		obj  Object
		want string
	}{
		{
			"pod",
			&Pod{ObjectMeta: ObjectMeta{Name: "web", Namespace: "default"}, Containers: []string{"app", "sidecar"}},
			"Pod default/web (2 containers)",
		},
		{
			"configmap",
			&ConfigMap{ObjectMeta: ObjectMeta{Name: "app-config", Namespace: "default"}, Data: map[string]string{"a": "1", "b": "2", "c": "3"}},
			"ConfigMap default/app-config (3 keys)",
		},
		{
			"unknown type falls back to default",
			&Service{ObjectMeta: ObjectMeta{Name: "web", Namespace: "prod"}},
			"Service prod/web",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Describe(tt.obj); got != tt.want {
				t.Errorf("Describe() = %q, want %q", got, tt.want)
			}
		})
	}
}
