package objects

import (
	"reflect"
	"testing"
)

func pod(ns, name string) *Pod {
	return &Pod{ObjectMeta: ObjectMeta{Name: name, Namespace: ns}}
}

func TestStoreCreateGet(t *testing.T) {
	s := NewStore()
	if err := s.Create(pod("default", "web")); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	obj, err := s.Get("Pod", "default", "web")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	p, ok := obj.(*Pod) // type assertion: Object -> *Pod
	if !ok {
		t.Fatalf("Get() returned %T, want *Pod", obj)
	}
	if p.Name != "web" {
		t.Errorf("Get() name = %q, want %q", p.Name, "web")
	}
}

func TestStoreCreateDuplicate(t *testing.T) {
	s := NewStore()
	original := pod("default", "web")
	original.Containers = []string{"app"}
	if err := s.Create(original); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := s.Create(pod("default", "web"))
	if !IsAlreadyExists(err) {
		t.Fatalf("second Create() error = %v, want AlreadyExists", err)
	}

	obj, _ := s.Get("Pod", "default", "web")
	if p, ok := obj.(*Pod); !ok || len(p.Containers) != 1 {
		t.Errorf("duplicate Create() overwrote the original object")
	}
}

func TestStoreSameNameDifferentKindOrNamespace(t *testing.T) {
	s := NewStore()
	objs := []Object{
		pod("default", "web"),
		pod("staging", "web"),
		&ConfigMap{ObjectMeta: ObjectMeta{Name: "web", Namespace: "default"}},
	}
	for _, o := range objs {
		if err := s.Create(o); err != nil {
			t.Fatalf("Create(%s %s/%s) error = %v", o.GetKind(), o.GetNamespace(), o.GetName(), err)
		}
	}
	obj, err := s.Get("ConfigMap", "default", "web")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if _, ok := obj.(*ConfigMap); !ok {
		t.Errorf("Get(ConfigMap) returned %T", obj)
	}
}

func TestStoreGetMissing(t *testing.T) {
	s := NewStore()
	obj, err := s.Get("Pod", "default", "nope")
	if !IsNotFound(err) {
		t.Errorf("Get() error = %v, want NotFound", err)
	}
	if obj != nil {
		t.Errorf("Get() obj = %v, want nil", obj)
	}
}

func TestStoreDelete(t *testing.T) {
	s := NewStore()
	_ = s.Create(pod("default", "web"))

	if err := s.Delete("Pod", "default", "web"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := s.Get("Pod", "default", "web"); !IsNotFound(err) {
		t.Errorf("Get() after Delete() error = %v, want NotFound", err)
	}
	if err := s.Delete("Pod", "default", "web"); !IsNotFound(err) {
		t.Errorf("second Delete() error = %v, want NotFound", err)
	}
}

func TestStoreList(t *testing.T) {
	s := NewStore()
	for _, o := range []Object{
		pod("staging", "b"),
		pod("default", "c"),
		pod("default", "a"),
		pod("staging", "a"),
		&ConfigMap{ObjectMeta: ObjectMeta{Name: "cfg", Namespace: "default"}},
	} {
		_ = s.Create(o)
	}

	ids := func(objs []Object) []string {
		var out []string
		for _, o := range objs {
			out = append(out, o.GetNamespace()+"/"+o.GetName())
		}
		return out
	}

	tests := []struct {
		name, kind, namespace string
		want                  []string
	}{
		{"one namespace", "Pod", "default", []string{"default/a", "default/c"}},
		{"all namespaces", "Pod", "", []string{"default/a", "default/c", "staging/a", "staging/b"}},
		{"other kind", "ConfigMap", "", []string{"default/cfg"}},
		{"nothing", "Secret", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run several times: random map order should never change the result.
			for i := 0; i < 5; i++ {
				if got := ids(s.List(tt.kind, tt.namespace)); !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("List(%q, %q) = %v, want %v", tt.kind, tt.namespace, got, tt.want)
				}
			}
		})
	}
}
