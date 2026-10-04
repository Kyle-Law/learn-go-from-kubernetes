package objects

import "sort"

// Store is a toy in-memory API server. It stores any Object, keyed by
// kind + namespace + name (e.g. "Pod/default/web").
type Store struct {
	objects map[string]Object
}

func NewStore() *Store {
	return &Store{objects: make(map[string]Object)}
}

func key(kind, namespace, name string) string {
	return kind + "/" + namespace + "/" + name
}

// Create stores obj. If an object with the same kind/namespace/name already
// exists, return an AlreadyExists error and leave the existing one in place.
func (s *Store) Create(obj Object) error {
	// TODO
	// PseudoCode
	// Check if the obj exists in s.object[key]
	// if exists, then return AlreadyExists error
	k := key(obj.GetKind(), obj.GetNamespace(), obj.GetName())
	if _, ok := s.objects[k]; ok {
		return NewAlreadyExists(obj.GetKind(), obj.GetName())
	}
	s.objects[k] = obj
	return nil
}

// Get returns the object, or (nil, NotFound error) if it doesn't exist.
// The caller gets an Object back and uses a type assertion to get the
// concrete type: pod := obj.(*Pod)
func (s *Store) Get(kind, namespace, name string) (Object, error) {
	// TODO
	// PseudoCode, if there's the key inside s.objects, if return the object; else return error
	k := key(kind, namespace, name)
	if obj, ok := s.objects[k]; ok {
		return obj, nil
	}
	return nil, NewNotFound(kind, name)
}

// Delete removes the object, or returns a NotFound error if it doesn't exist.
func (s *Store) Delete(kind, namespace, name string) error {
	// TODO
	k := key(kind, namespace, name)
	if _, ok := s.objects[k]; ok {
		delete(s.objects, k)
		return nil
	}
	return NewNotFound(kind, name)
}

func IsSameNameSpace(obj Object, namespace string) bool {
	return namespace == "" || obj.GetNamespace() == namespace
}

// List returns all objects of kind in namespace ("" = all namespaces),
// sorted by namespace, then name. Map iteration order in Go is random,
// so you MUST sort. Hint: sort.Slice or slices.SortFunc.
func (s *Store) List(kind, namespace string) []Object {
	var out []Object
	// FILTER & APPEND
	for _, obj := range s.objects {
		if obj.GetKind() == kind && IsSameNameSpace(obj, namespace) {
			out = append(out, obj)
		}
	}
	// SORT
	sort.Slice(out, func(i, j int) bool {
		if out[i].GetNamespace() != out[j].GetNamespace() {
			return out[i].GetNamespace() < out[j].GetNamespace()
		}
		return out[i].GetName() < out[j].GetName()
	})

	return out
}
