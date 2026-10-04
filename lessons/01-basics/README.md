# 01 — Basics: modelling a Pod

Goal: implement the `TODO`s in `pod.go` until `go test ./lessons/01-basics/` passes.

## Concepts you need

### Structs = resource schemas
```go
type Container struct {
    Name  string
    Image string
    Ready bool
}
```
Capitalized names are **exported** (public outside the package). Lowercase = private.
That's the whole visibility system.

### Zero values = "field not set"
Every type has a zero value: `""`, `0`, `false`, `nil`. A `Pod{}` is valid with all fields
zeroed — like applying a manifest with fields omitted. No constructors required.

### Slices = lists (`spec.containers`)
```go
cs := []Container{{Name: "app"}, {Name: "sidecar"}}
cs = append(cs, Container{Name: "init"})   // append returns the new slice — always reassign
for i, c := range cs { ... }               // c is a COPY of the element
for _, c := range cs { ... }               // _ discards the index
```

### Maps = labels
```go
labels := map[string]string{"app": "web"}
v, ok := labels["tier"]   // "comma ok": ok is false if key missing, v is ""
```
⚠️ A **nil map** can be read but not written: `var m map[string]string; m["a"] = "b"` panics.
Initialize with `make(map[string]string)` or a literal first.

### Methods and receivers
```go
func (p Pod) Ready() bool      { ... }  // value receiver: gets a copy, can't modify p
func (p *Pod) SetLabel(k, v string) { ... }  // pointer receiver: modifies the real Pod
```
Rule of thumb: if a method mutates, use a pointer receiver.

### Named types + constants = enums
```go
type Phase string
const (
    PodRunning Phase = "Running"
)
```
Same pattern as `corev1.PodPhase` in the real Kubernetes code.

### Formatting strings
`fmt.Sprintf("%d/%d", ready, total)` — `%d` int, `%s` string, `%v` anything, `%+v` struct with field names.

## Exercises (in `pod.go`)
1. `Pod.Ready()` — what makes a pod "Ready"?
2. `Pod.ReadyString()` — the `READY` column of `kubectl get pods`, e.g. `1/2`.
3. `SelectorMatches()` — equality-based label selector (`-l app=web,tier=frontend`).
4. `FilterPods()` — `kubectl get pods -n <ns> -l <selector>`.
5. `CountByPhase()` — a quick status summary.
6. `Pod.SetLabel()` — careful with that nil map.

## Bonus
- Run `go vet ./...` and `gofmt -l .` — Go's built-in linters. Most editors run `gofmt` on save.
- Read `pod_test.go`. The **table-driven test** style there is how nearly all Go (and Kubernetes) tests are written.
- Look at the real thing: [`labels.SelectorFromSet`](https://github.com/kubernetes/apimachinery/blob/master/pkg/labels/selector.go).
