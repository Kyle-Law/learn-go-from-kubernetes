# 03 — Encoding: parsing a Deployment manifest

Goal: make `go test ./lessons/03-encoding/` pass.
1. **`types.go`**: add struct tags so the types read and write real Kubernetes YAML and JSON.
2. **`manifest.go`**: parse, encode, default, edit, and validate a Deployment.

`testdata/deployment.yaml` is a normal manifest. You could `kubectl apply` it.

## Concepts you need

### Struct tags
A tag is a string literal after a field. Libraries read it at runtime with reflection:
```go
type ObjectMeta struct {
    Name      string            `json:"name"`
    Namespace string            `json:"namespace,omitempty"`
    Labels    map[string]string `json:"labels,omitempty"`
    internal  string            // unexported: never encoded or decoded
    Secret    string            `json:"-"` // exported, but explicitly skipped
}
```
- `json:"name"` sets the key name. Without a tag, the key is the Go field name (`"Name"`).
- `,omitempty` skips the field when it's the zero value: `""`, `0`, `false`, `nil`, or an
  empty map or slice. ⚠️ It does **not** skip an empty *struct*. Use a pointer if a struct
  must disappear when unset.
- Only **exported** (capitalized) fields are encoded. A lowercase field is silently ignored.

Open [`k8s.io/api/apps/v1/types.go`](https://github.com/kubernetes/api/blob/master/apps/v1/types.go)
and you'll see these same tags on every field. They're also why `kubectl explain` and the YAML
you write use camelCase.

### encoding/json
```go
data, err := json.Marshal(v)          // Go value -> []byte
err := json.Unmarshal(data, &v)       // []byte -> Go value (pass a POINTER)
```
Two traps:
- **Unmarshal matches keys case-insensitively.** `{"METADATA": ...}` still fills `Metadata`.
  So parsing can look fine before your tags are right. Marshal shows the truth.
- **Unknown keys are silently dropped.** Write `replcas: 3` and nothing complains; you just
  get the default. (This is why `kubectl apply` now has server-side field validation.)

### YAML in Kubernetes = YAML → JSON → struct
`sigs.k8s.io/yaml` converts YAML to JSON and then calls `encoding/json`. That's why Kubernetes
types only need `json` tags, and why JSON input works too: JSON is valid YAML.

### Pointers for optional fields
```go
Replicas *int32 `json:"replicas,omitempty"`
```
With a plain `int32`, `replicas: 0` (scaled to zero) and "not specified" (defaults to 1)
would look the same. A pointer gives three states: `nil` (unset), `&0`, and `&3`. All the
`*int32` and `*bool` fields in the Kubernetes API exist for this reason.

You can't write `&3` in Go, which is why the `Ptr` helper exists:
```go
d.Spec.Replicas = Ptr(int32(3))
if d.Spec.Replicas != nil { n := *d.Spec.Replicas }   // dereference with *
```

### A first taste of generics
```go
func Ptr[T any](v T) *T { ... }
```
`[T any]` is a type parameter: the same function works for any type, and Go infers `T` from
the argument. Kubernetes uses this in `k8s.io/utils/ptr.To`. You'll see generics more in
lesson 8. For now, this is all you need.

### Modifying slice elements
```go
for _, c := range containers { c.Image = "x" }        // ❌ c is a copy; nothing changes
for i := range containers { containers[i].Image = "x" } // ✅ index into the slice
```

### Combining errors
```go
var errs []error
errs = append(errs, ErrMissingName)
errs = append(errs, fmt.Errorf("container %q: %w", name, ErrMissingImage))
return errors.Join(errs...)   // nil if errs is empty; errors.Is works on each one
```
`kubectl apply` shows every validation problem at once instead of stopping at the first.
`errors.Join` gives you the same behavior.

### `[]byte` and reading files
`os.ReadFile(path)` returns `([]byte, error)`. Encoders work on `[]byte`, and you convert with
`string(b)` or `[]byte(s)`. Tests read files from `testdata/`, a directory name the Go tool
treats specially and ignores when it builds packages.

## Exercises
1. Tags in `types.go`. Check with `go test -run 'TestToJSON|TestRoundTrip' ./lessons/03-encoding/`.
2. `ParseDeployment`: wrap the parse error, check apiVersion and kind.
3. `ToJSON`, `Ptr`, `ReplicasOrDefault`.
4. `SetImage`: watch out for the range copy.
5. `Validate` with `errors.Join`.

## Bonus
- Build a tiny CLI: `go run ./lessons/03-encoding/cmd/validate testdata/deployment.yaml`
  that prints `OK` or each validation error. (Hint: `os.Args`, `os.Exit(1)`.)
- Try `kubectl create deployment web --image=nginx --dry-run=client -o yaml | <your CLI>`.
- What happens if you remove `,omitempty` from `Status` and marshal a Deployment without
  a status? Why does `kubectl get -o yaml` show `status: {}` for some objects?
