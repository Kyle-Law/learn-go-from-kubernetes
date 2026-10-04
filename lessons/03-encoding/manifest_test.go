package manifests

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func readFile(t *testing.T, path string) []byte {
	t.Helper() // failures are reported at the caller's line, not here
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return data
}

func parseTestdata(t *testing.T) *Deployment {
	t.Helper()
	d, err := ParseDeployment(readFile(t, "testdata/deployment.yaml"))
	if err != nil {
		t.Fatalf("ParseDeployment() error = %v", err)
	}
	if d == nil {
		t.Fatal("ParseDeployment() returned nil Deployment")
	}
	return d
}

func TestParseDeployment(t *testing.T) {
	d := parseTestdata(t)

	if d.Kind != "Deployment" || d.APIVersion != "apps/v1" {
		t.Errorf("TypeMeta = %+v", d.TypeMeta)
	}
	if d.Metadata.Name != "web" || d.Metadata.Namespace != "default" {
		t.Errorf("Metadata = %+v", d.Metadata)
	}
	if d.Spec.Replicas == nil || *d.Spec.Replicas != 3 {
		t.Errorf("Spec.Replicas = %v, want pointer to 3", d.Spec.Replicas)
	}
	if got := d.Spec.Template.Metadata.Labels["tier"]; got != "frontend" {
		t.Errorf("template label tier = %q, want frontend", got)
	}
	cs := d.Spec.Template.Spec.Containers
	if len(cs) != 2 {
		t.Fatalf("got %d containers, want 2", len(cs))
	}
	if cs[0].Image != "nginx:1.27" || len(cs[0].Ports) != 1 || cs[0].Ports[0].ContainerPort != 80 {
		t.Errorf("containers[0] = %+v", cs[0])
	}
	if len(cs[0].Env) != 1 || cs[0].Env[0].Value != "info" {
		t.Errorf("containers[0].Env = %+v", cs[0].Env)
	}
	if d.Status != nil {
		t.Errorf("Status = %+v, want nil (not in the manifest)", d.Status)
	}
}

func TestParseDeploymentFromJSON(t *testing.T) {
	// JSON is valid YAML, so the same parser handles `kubectl get -o json` output.
	d, err := ParseDeployment([]byte(`{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"api"}}`))
	if err != nil {
		t.Fatalf("ParseDeployment() error = %v", err)
	}
	if d.Metadata.Name != "api" {
		t.Errorf("name = %q, want api", d.Metadata.Name)
	}
}

func TestParseRejectsOtherKinds(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantText string
	}{
		{"service", readFile(t, "testdata/service.yaml"), "v1/Service"},
		{"old apiVersion", []byte("apiVersion: extensions/v1beta1\nkind: Deployment\n"), "extensions/v1beta1/Deployment"},
		{"empty document", []byte(""), "unsupported kind"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ParseDeployment(tt.input)
			if !errors.Is(err, ErrUnsupportedKind) {
				t.Fatalf("error = %v, want one wrapping ErrUnsupportedKind", err)
			}
			if !strings.Contains(err.Error(), tt.wantText) {
				t.Errorf("error %q should mention %q", err, tt.wantText)
			}
			if d != nil {
				t.Errorf("Deployment = %+v, want nil on error", d)
			}
		})
	}
}

func TestParseInvalidYAML(t *testing.T) {
	_, err := ParseDeployment([]byte("metadata: [unclosed"))
	if err == nil {
		t.Fatal("expected an error for invalid YAML")
	}
	if !strings.HasPrefix(err.Error(), "parsing deployment: ") {
		t.Errorf("error = %q, want prefix %q", err, "parsing deployment: ")
	}
	if errors.Unwrap(err) == nil {
		t.Errorf("error %q doesn't wrap the original error (use %%w)", err)
	}
	if errors.Is(err, ErrUnsupportedKind) {
		t.Errorf("invalid YAML should not be reported as ErrUnsupportedKind")
	}
}

func TestPtr(t *testing.T) {
	p := Ptr(int32(3))
	if p == nil || *p != 3 {
		t.Fatalf("Ptr(int32(3)) = %v", p)
	}
	s := Ptr("hello")
	if s == nil || *s != "hello" {
		t.Fatalf(`Ptr("hello") = %v`, s)
	}
	if Ptr(1) == Ptr(1) {
		t.Errorf("each call should return a new pointer")
	}
}

func TestReplicasOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		replicas *int32
		want     int32
	}{
		{"unset defaults to 1", nil, 1},
		{"explicit zero stays zero", Ptr(int32(0)), 0},
		{"explicit value", Ptr(int32(5)), 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deployment{Spec: DeploymentSpec{Replicas: tt.replicas}}
			if got := d.ReplicasOrDefault(); got != tt.want {
				t.Errorf("ReplicasOrDefault() = %d, want %d", got, tt.want)
			}
		})
	}

	t.Run("replicas: 0 in YAML", func(t *testing.T) {
		d, err := ParseDeployment([]byte("apiVersion: apps/v1\nkind: Deployment\nspec:\n  replicas: 0\n"))
		if err != nil {
			t.Fatalf("ParseDeployment() error = %v", err)
		}
		if got := d.ReplicasOrDefault(); got != 0 {
			t.Errorf("ReplicasOrDefault() = %d, want 0 (scaled to zero, not unset)", got)
		}
	})
}

func minimalDeployment() *Deployment {
	return &Deployment{
		TypeMeta: TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		Metadata: ObjectMeta{Name: "web"},
		Spec: DeploymentSpec{
			Selector: LabelSelector{MatchLabels: map[string]string{"app": "web"}},
			Template: PodTemplateSpec{
				Metadata: ObjectMeta{Labels: map[string]string{"app": "web"}},
				Spec: PodSpec{Containers: []Container{
					{Name: "app", Image: "nginx", Ports: []ContainerPort{{ContainerPort: 80}}},
				}},
			},
		},
	}
}

func TestToJSON(t *testing.T) {
	t.Run("tags and omitempty", func(t *testing.T) {
		got, err := ToJSON(minimalDeployment())
		if err != nil {
			t.Fatalf("ToJSON() error = %v", err)
		}
		want := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"web"},` +
			`"spec":{"selector":{"matchLabels":{"app":"web"}},` +
			`"template":{"metadata":{"labels":{"app":"web"}},` +
			`"spec":{"containers":[{"name":"app","image":"nginx","ports":[{"containerPort":80}]}]}}}}`
		if string(got) != want {
			t.Errorf("ToJSON()\n got: %s\nwant: %s", got, want)
		}
	})

	t.Run("replicas set to zero is kept", func(t *testing.T) {
		d := minimalDeployment()
		d.Spec.Replicas = Ptr(int32(0))
		got, _ := ToJSON(d)
		if !strings.Contains(string(got), `"spec":{"replicas":0,`) {
			t.Errorf("ToJSON() = %s\nwant it to contain \"replicas\":0", got)
		}
	})

	t.Run("status when set", func(t *testing.T) {
		d := minimalDeployment()
		d.Status = &DeploymentStatus{ReadyReplicas: 2}
		got, _ := ToJSON(d)
		if !strings.HasSuffix(string(got), `,"status":{"readyReplicas":2}}`) {
			t.Errorf("ToJSON() = %s\nwant it to end with \"status\":{\"readyReplicas\":2}", got)
		}
	})
}

func TestRoundTrip(t *testing.T) {
	original := parseTestdata(t)
	data, err := ToJSON(original)
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}
	back, err := ParseDeployment(data)
	if err != nil {
		t.Fatalf("ParseDeployment(ToJSON()) error = %v\nJSON: %s", err, data)
	}
	if !reflect.DeepEqual(original, back) {
		t.Errorf("round trip changed the object\n got: %+v\nwant: %+v", back, original)
	}
}

func TestSetImage(t *testing.T) {
	d := parseTestdata(t)
	if err := SetImage(d, "sidecar", "busybox:1.37"); err != nil {
		t.Fatalf("SetImage() error = %v", err)
	}
	cs := d.Spec.Template.Spec.Containers
	if cs[1].Image != "busybox:1.37" {
		t.Errorf("sidecar image = %q, want busybox:1.37 (did you modify a copy?)", cs[1].Image)
	}
	if cs[0].Image != "nginx:1.27" {
		t.Errorf("app image changed to %q", cs[0].Image)
	}

	err := SetImage(d, "nope", "x")
	if !errors.Is(err, ErrContainerNotFound) {
		t.Errorf("SetImage(missing) error = %v, want ErrContainerNotFound", err)
	}
}

func TestValidate(t *testing.T) {
	if err := Validate(parseTestdata(t)); err != nil {
		t.Fatalf("Validate(valid) = %v, want nil", err)
	}

	all := []error{ErrMissingName, ErrNoContainers, ErrMissingImage, ErrSelectorMismatch}
	tests := []struct {
		name   string
		mutate func(d *Deployment)
		want   []error
	}{
		{"missing name", func(d *Deployment) { d.Metadata.Name = "" }, []error{ErrMissingName}},
		{"no containers", func(d *Deployment) { d.Spec.Template.Spec.Containers = nil }, []error{ErrNoContainers}},
		{"missing image", func(d *Deployment) { d.Spec.Template.Spec.Containers[0].Image = "" }, []error{ErrMissingImage}},
		{"empty selector", func(d *Deployment) { d.Spec.Selector.MatchLabels = nil }, []error{ErrSelectorMismatch}},
		{"selector doesn't match", func(d *Deployment) { d.Spec.Selector.MatchLabels["app"] = "api" }, []error{ErrSelectorMismatch}},
		{
			"several problems at once",
			func(d *Deployment) {
				d.Metadata.Name = ""
				d.Spec.Template.Spec.Containers[1].Image = ""
				d.Spec.Template.Metadata.Labels = nil
			},
			[]error{ErrMissingName, ErrMissingImage, ErrSelectorMismatch},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := parseTestdata(t)
			tt.mutate(d)
			err := Validate(d)
			for _, target := range all {
				want := false
				for _, w := range tt.want {
					if w == target {
						want = true
					}
				}
				if got := errors.Is(err, target); got != want {
					t.Errorf("errors.Is(err, %q) = %v, want %v\nerr: %v", target, got, want, err)
				}
			}
		})
	}

	t.Run("image error names the container", func(t *testing.T) {
		d := parseTestdata(t)
		d.Spec.Template.Spec.Containers[1].Image = ""
		if err := Validate(d); err == nil || !strings.Contains(err.Error(), `container "sidecar": image is required`) {
			t.Errorf("Validate() = %v, want it to mention container \"sidecar\"", err)
		}
	})
}
