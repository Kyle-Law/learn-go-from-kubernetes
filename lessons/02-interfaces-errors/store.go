package objects

// Store is a toy in-memory API server. It stores any Object, keyed by
// kind + namespace + name (e.g. "Pod/default/web").
type Store struct {
	objects map[string]Object
}

func NewStore() *Store {
	return &Store{objects: make(map[string]Object)}
}

// Create stores obj. If an object with the same kind/namespace/name already
// exists, return an AlreadyExists error and leave the existing one in place.
func (s *Store) Create(obj Object) error {
	// TODO
	return nil
}

// Get returns the object, or (nil, NotFound error) if it doesn't exist.
// The caller gets an Object back and uses a type assertion to get the
// concrete type: pod := obj.(*Pod)
func (s *Store) Get(kind, namespace, name string) (Object, error) {
	// TODO
	return nil, nil
}

// Delete removes the object, or returns a NotFound error if it doesn't exist.
func (s *Store) Delete(kind, namespace, name string) error {
	// TODO
	return nil
}

// List returns all objects of kind in namespace ("" = all namespaces),
// sorted by namespace, then name. Map iteration order in Go is random,
// so you MUST sort. Hint: sort.Slice or slices.SortFunc.
func (s *Store) List(kind, namespace string) []Object {
	// TODO
	return nil
}
