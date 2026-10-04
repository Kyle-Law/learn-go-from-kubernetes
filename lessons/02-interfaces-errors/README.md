# 02 — Interfaces & errors: a toy API server

Goal: implement the `TODO`s in `object.go`, `errors.go`, and `store.go` until
`go test ./lessons/02-interfaces-errors/` passes. Suggested order is the same:
objects → errors → store.

## Concepts you need

### Interfaces are implicit
```go
type Object interface {
    GetName() string
    GetKind() string
}
```
Any type with those methods **is** an `Object` — no `implements` keyword. That's how
client-go can take any resource (`Pod`, `Deployment`, your CRD) as a `client.Object`.
See the `Service` type in `object_test.go`: it's defined in the test file and still
works with your `Describe`.

Keep interfaces small. The most-used interfaces in Go have one method:
`error` (`Error() string`), `fmt.Stringer` (`String() string`), `io.Reader` (`Read`).

### Embedding = metadata on every resource
```go
type Pod struct {
    ObjectMeta          // no field name -> embedded
    Containers []string
}
p.Name          // promoted field (really p.ObjectMeta.Name)
p.GetName()     // promoted method
```
This is exactly how every real Kubernetes type gets `metav1.ObjectMeta` and
`metav1.TypeMeta`. It's composition, not inheritance: a `*Pod` is not an `*ObjectMeta`.

### Type assertions and type switches
Going from an interface back to a concrete type:
```go
p, ok := obj.(*Pod)      // ok == false if obj isn't a *Pod (no panic)
p := obj.(*Pod)          // panics if wrong — avoid unless you're certain

switch o := obj.(type) { // o has the concrete type inside each case
case *Pod:
    return len(o.Containers)
case *ConfigMap:
    return len(o.Data)
default:
    return 0
}
```
Informer event handlers give you `interface{}` and you do exactly this.

### Errors are values
There are no exceptions. Functions return an `error` as the last value, and you check it:
```go
obj, err := store.Get("Pod", "default", "web")
if err != nil {
    return fmt.Errorf("getting pod: %w", err)   // add context and wrap
}
```
- `error` is just an interface: `type error interface { Error() string }`.
  Any type with that method is an error — that's how `*StatusError` works.
- `%w` **wraps** the error (keeps the original inside). `%v` only copies the text.
- `errors.Is(err, target)` checks whether a specific error value is in the chain (good for
  sentinels like `io.EOF`).
- `errors.As(err, &target)` checks whether an error of a given **type** is in the chain and
  extracts it:
  ```go
  var se *StatusError
  if errors.As(err, &se) {
      fmt.Println(se.Code)
  }
  ```
- Never compare error **text** (`err.Error() == "..."`). One of the tests checks that you don't.

### Multiple return values and "comma ok"
`(Object, error)` returns are everywhere. Convention: when `err != nil`, the other values are
zero values (`nil` here), and callers shouldn't use them.

## Exercises
1. `object.go` — `ObjectMeta` getters, `GetKind` for `Pod`/`ConfigMap`, `Describe` with a type switch.
2. `errors.go` — `StatusError`, constructors, and `IsNotFound` / `IsAlreadyExists` that work on wrapped errors.
3. `store.go` — `Create` / `Get` / `Delete` / `List` for the toy API server.

## Bonus
- Change `var _ Object = &Pod{}` in `object_test.go` to `var _ Object = Pod{}` (no `&`) and
  read the compile error. Why? (Hint: method sets. The methods have pointer receivers.)
- The **nil interface trap**: what does this print, and why?
  ```go
  func find() error {
      var se *StatusError = nil
      return se
  }
  fmt.Println(find() == nil)
  ```
  This is why the constructors return `*StatusError`, but functions like `Store.Get` should
  `return nil, NewNotFound(...)` directly, and never return a typed nil pointer as an `error`.
- Read the real thing: [`apimachinery/pkg/api/errors/errors.go`](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/errors/errors.go)
  — look for `IsNotFound` and `ReasonForError`.
