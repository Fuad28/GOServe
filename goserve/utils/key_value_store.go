package utils

// KeyValueStore is a type that allows proper storing and retrieval of data.
// It's used accross the project for headers, query & path parameters.
// The implementaiton is generic and exposed to be used outside the project.
type KeyValueStore[K comparable, T any] struct {

	// data is a map that holds data collection
	data map[K]T
}

// NewKeyValueStore initializes and returns a pointer to a new KeyValueStore instance.
func NewKeyValueStore[K comparable, T comparable]() *KeyValueStore[K, T] {
	return &KeyValueStore[K, T]{
		data: make(map[K]T),
	}
}

// Set adds or updates a key-value pair.
func (kv *KeyValueStore[K, T]) Set(key K, value T) {
	kv.data[key] = value
}

// Get retrieves the value associated with a key.
func (kv *KeyValueStore[K, T]) Get(key K) (T, bool) {
	value, exists := kv.data[key]
	return value, exists
}

// Delete removes a key-value pair.
func (kv *KeyValueStore[K, T]) Delete(key K) {
	delete(kv.data, key)
}

// Has checks if a key exists in the store.
func (kv *KeyValueStore[K, T]) Has(key K) bool {
	_, exists := kv.data[key]
	return exists
}

// GetAll returns all key-value pairs.
func (kv *KeyValueStore[K, T]) GetAll() map[K]T {
	return kv.data
}
