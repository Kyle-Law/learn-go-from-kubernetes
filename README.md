# Go from Kubernetes

Learn Go by rebuilding the parts of Kubernetes you already know.
Each lesson has a `README.md` (concepts), a `.go` file with `TODO`s, and tests.
You're done with a lesson when its tests pass:

```sh
go test ./lessons/01-basics/...        # run one lesson
go test -v -run TestSelector ./lessons/01-basics/   # run one test, verbose
```

## Layout

- `lessons/NN-*/` — blank exercises. Start here.
- `attempts/NN-*/` — my completed attempts, plus a `go-mistakes-quiz.html` per lesson
  (open it in a browser) covering the tricky parts and the mistakes I made.
  Spoilers: finish the lesson before looking.

## Roadmap

| #  | Lesson                         | Go concepts                                              | The K8s thing you build                          |
|----|--------------------------------|----------------------------------------------------------|--------------------------------------------------|
| 01 | Basics                         | types, structs, methods, slices, maps, pointers          | Pod model, label selectors, `kubectl get` READY  |
| 02 | Interfaces & errors            | interfaces, embedding, `errors.Is/As`, wrapping           | `runtime.Object`, `IsNotFound`-style API errors   |
| 03 | Encoding                       | struct tags, `encoding/json`, `omitempty`, YAML          | Parse a Deployment manifest                      |
| 04 | Concurrency                    | goroutines, channels, `select`, `sync`, `context`         | A work queue + reconcile loop                    |
| 05 | HTTP services                  | `net/http`, signals, graceful shutdown                   | App with `/healthz` `/readyz`, SIGTERM handling, run in kind |
| 06 | client-go                      | modules, real-world API usage                            | List/watch Pods against your kind cluster        |
| 07 | Informers & controllers        | putting it all together                                  | A controller that labels/annotates Pods          |
| 08 | Operators                      | kubebuilder / controller-runtime, generics               | A CRD + operator                                 |

## Mental map: K8s → Go

| You know (K8s)                      | Go equivalent                                         |
|-------------------------------------|-------------------------------------------------------|
| A resource schema (`PodSpec`)       | a `struct`                                            |
| `metadata.labels`                   | `map[string]string`                                   |
| `spec.containers[]`                 | a slice: `[]Container`                                |
| Field omitted in YAML → default     | Zero values (`""`, `0`, `false`, `nil`)               |
| `kind: Pod` / duck typing on `kind` | interfaces (implicit — no `implements` keyword)       |
| Controller reconcile loop           | `for` + channels + `context.Context`                  |
| `kubectl get pod -o json`           | `json.Marshal` driven by struct tags                  |
| Liveness/readiness probes, SIGTERM  | `net/http` handlers, `os/signal`                      |

## Recommended side reading
- [A Tour of Go](https://go.dev/tour/) — do it alongside lessons 01–04
- [Effective Go](https://go.dev/doc/effective_go)
- The `k8s.io/api/core/v1` types — once you've done lesson 03, read `types.go` there; it'll suddenly make sense
