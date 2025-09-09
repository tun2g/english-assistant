package patterns

import (
	"fmt"
	"sync"
)

// ConcurrentMap provides a thread-safe map implementation
type ConcurrentMap[K comparable, V any] struct {
	shards   []*shard[K, V]
	shardMask uint32
}

// shard represents a single shard of the concurrent map
type shard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

const defaultShardCount = 32

// NewConcurrentMap creates a new concurrent map with default shard count
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return NewConcurrentMapWithShards[K, V](defaultShardCount)
}

// NewConcurrentMapWithShards creates a new concurrent map with specified shard count
func NewConcurrentMapWithShards[K comparable, V any](shardCount uint32) *ConcurrentMap[K, V] {
	// Ensure shard count is power of 2 for efficient modulo
	shardCount = nextPowerOf2(shardCount)
	
	cm := &ConcurrentMap[K, V]{
		shards:   make([]*shard[K, V], shardCount),
		shardMask: shardCount - 1,
	}
	
	for i := uint32(0); i < shardCount; i++ {
		cm.shards[i] = &shard[K, V]{
			items: make(map[K]V),
		}
	}
	
	return cm
}

// getShard returns the shard for the given key
func (cm *ConcurrentMap[K, V]) getShard(key K) *shard[K, V] {
	return cm.shards[hash(key)&cm.shardMask]
}

// Set sets a key-value pair
func (cm *ConcurrentMap[K, V]) Set(key K, value V) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.items[key] = value
}

// Get retrieves a value by key
func (cm *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shard := cm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	value, ok := shard.items[key]
	return value, ok
}

// Delete removes a key-value pair
func (cm *ConcurrentMap[K, V]) Delete(key K) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

// Has checks if a key exists
func (cm *ConcurrentMap[K, V]) Has(key K) bool {
	_, ok := cm.Get(key)
	return ok
}

// Size returns the total number of items in the map
func (cm *ConcurrentMap[K, V]) Size() int {
	size := 0
	for _, shard := range cm.shards {
		shard.mu.RLock()
		size += len(shard.items)
		shard.mu.RUnlock()
	}
	return size
}

// Keys returns all keys in the map
func (cm *ConcurrentMap[K, V]) Keys() []K {
	keys := make([]K, 0, cm.Size())
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for key := range shard.items {
			keys = append(keys, key)
		}
		shard.mu.RUnlock()
	}
	return keys
}

// Values returns all values in the map
func (cm *ConcurrentMap[K, V]) Values() []V {
	values := make([]V, 0, cm.Size())
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for _, value := range shard.items {
			values = append(values, value)
		}
		shard.mu.RUnlock()
	}
	return values
}

// Items returns all key-value pairs
func (cm *ConcurrentMap[K, V]) Items() map[K]V {
	items := make(map[K]V)
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for key, value := range shard.items {
			items[key] = value
		}
		shard.mu.RUnlock()
	}
	return items
}

// Clear removes all items from the map
func (cm *ConcurrentMap[K, V]) Clear() {
	for _, shard := range cm.shards {
		shard.mu.Lock()
		shard.items = make(map[K]V)
		shard.mu.Unlock()
	}
}

// GetOrSet gets a value or sets it if it doesn't exist
func (cm *ConcurrentMap[K, V]) GetOrSet(key K, defaultValue V) (V, bool) {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	
	if value, ok := shard.items[key]; ok {
		return value, false // Value existed
	}
	
	shard.items[key] = defaultValue
	return defaultValue, true // Value was set
}

// GetOrCompute gets a value or computes it if it doesn't exist
func (cm *ConcurrentMap[K, V]) GetOrCompute(key K, computeFn func() V) V {
	shard := cm.getShard(key)
	
	// First try with read lock
	shard.mu.RLock()
	if value, ok := shard.items[key]; ok {
		shard.mu.RUnlock()
		return value
	}
	shard.mu.RUnlock()
	
	// Upgrade to write lock
	shard.mu.Lock()
	defer shard.mu.Unlock()
	
	// Double-check after acquiring write lock
	if value, ok := shard.items[key]; ok {
		return value
	}
	
	// Compute and set the value
	value := computeFn()
	shard.items[key] = value
	return value
}

// ForEach iterates over all key-value pairs
func (cm *ConcurrentMap[K, V]) ForEach(fn func(K, V) bool) {
	for _, shard := range cm.shards {
		shard.mu.RLock()
		for key, value := range shard.items {
			if !fn(key, value) {
				shard.mu.RUnlock()
				return
			}
		}
		shard.mu.RUnlock()
	}
}

// Filter returns a new map with items that match the predicate
func (cm *ConcurrentMap[K, V]) Filter(predicate func(K, V) bool) *ConcurrentMap[K, V] {
	result := NewConcurrentMapWithShards[K, V](uint32(len(cm.shards)))
	
	cm.ForEach(func(key K, value V) bool {
		if predicate(key, value) {
			result.Set(key, value)
		}
		return true
	})
	
	return result
}

// Update atomically updates a value if the key exists
func (cm *ConcurrentMap[K, V]) Update(key K, updateFn func(V) V) bool {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	
	if value, ok := shard.items[key]; ok {
		shard.items[key] = updateFn(value)
		return true
	}
	return false
}

// CompareAndSwap atomically compares and swaps a value
func (cm *ConcurrentMap[K, V]) CompareAndSwap(key K, oldValue, newValue V, compareFn func(V, V) bool) bool {
	shard := cm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	
	if currentValue, ok := shard.items[key]; ok && compareFn(currentValue, oldValue) {
		shard.items[key] = newValue
		return true
	}
	return false
}

// nextPowerOf2 returns the next power of 2 greater than or equal to n
func nextPowerOf2(n uint32) uint32 {
	if n == 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

// hash provides a simple hash function for any comparable type
func hash[K comparable](key K) uint32 {
	// This is a simple hash function - in production, you might want something more sophisticated
	h := uint32(0)
	data := []byte(fmt.Sprintf("%v", key))
	for _, b := range data {
		h = h*31 + uint32(b)
	}
	return h
}

