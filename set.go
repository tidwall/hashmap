package hashmap

type Set[K comparable] struct {
	base Map[K, struct{}]
}

// Insert an item
func (tr *Set[K]) Insert(key K) {
	tr.base.Set(key, struct{}{})
}

// Get a value for key
func (tr *Set[K]) Contains(key K) bool {
	_, ok := tr.base.Get(key)
	return ok
}

// Len returns the number of items in the tree
func (tr *Set[K]) Len() int {
	return tr.base.Len()
}

// Delete an item
func (tr *Set[K]) Delete(key K) {
	tr.base.Delete(key)
}

func (tr *Set[K]) Scan(iter func(key K) bool) {
	tr.base.Scan(func(key K, value struct{}) bool {
		return iter(key)
	})
}

// Keys returns all keys as a slice
func (tr *Set[K]) Keys() []K {
	return tr.base.Keys()
}
