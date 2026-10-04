package objects

// Object is a tiny version of runtime.Object / metav1.Object: anything the
// API server can store. There is no "implements" keyword in Go — any type
// with these three methods IS an Object.
type Object interface {
	GetName() string
	GetNamespace() string
	GetKind() string
}

// ObjectMeta mirrors metav1.ObjectMeta. Every resource embeds it.
type ObjectMeta struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

// GetName returns the object's name.
// Because Pod and ConfigMap embed ObjectMeta, they get this method for free.
func (m *ObjectMeta) GetName() string {
	// TODO
	return ""
}

// GetNamespace returns the object's namespace.
func (m *ObjectMeta) GetNamespace() string {
	// TODO
	return ""
}

type Pod struct {
	ObjectMeta // embedded: fields AND methods are promoted onto Pod
	Containers []string
}

// GetKind returns "Pod".
func (p *Pod) GetKind() string {
	// TODO
	return ""
}

type ConfigMap struct {
	ObjectMeta
	Data map[string]string
}

// GetKind returns "ConfigMap".
func (c *ConfigMap) GetKind() string {
	// TODO
	return ""
}

// Describe returns a one-line summary, like a very short `kubectl describe`.
// Use a type switch to handle each concrete type:
//
//	*Pod       -> "Pod default/web (2 containers)"
//	*ConfigMap -> "ConfigMap default/app-config (3 keys)"
//	anything else -> "<Kind> <namespace>/<name>", e.g. "Service default/web"
func Describe(obj Object) string {
	// TODO
	return ""
}
