package patterns_test

import (
	"sync"
	"testing"

	"app-backend/pkg/patterns"
)

func TestConcurrentMap(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()

		// Test Set and Get
		cm.Set("key1", 100)
		if value, exists := cm.Get("key1"); !exists || value != 100 {
			t.Errorf("Expected key1=100, got exists=%v, value=%d", exists, value)
		}

		// Test Has
		if !cm.Has("key1") {
			t.Error("Expected key1 to exist")
		}

		if cm.Has("nonexistent") {
			t.Error("Expected nonexistent key to not exist")
		}

		// Test Delete
		cm.Delete("key1")
		if cm.Has("key1") {
			t.Error("Expected key1 to be deleted")
		}
	})

	t.Run("get or set", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()

		// First call should set the value
		value, wasSet := cm.GetOrSet("key1", 42)
		if !wasSet || value != 42 {
			t.Errorf("Expected wasSet=true, value=42, got wasSet=%v, value=%d", wasSet, value)
		}

		// Second call should get existing value
		value, wasSet = cm.GetOrSet("key1", 99)
		if wasSet || value != 42 {
			t.Errorf("Expected wasSet=false, value=42, got wasSet=%v, value=%d", wasSet, value)
		}
	})

	t.Run("get or compute", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()

		computeCalled := false
		value := cm.GetOrCompute("key1", func() int {
			computeCalled = true
			return 123
		})

		if !computeCalled || value != 123 {
			t.Errorf("Expected compute to be called, value=123, got called=%v, value=%d", computeCalled, value)
		}

		// Second call should not compute
		computeCalled = false
		value = cm.GetOrCompute("key1", func() int {
			computeCalled = true
			return 456
		})

		if computeCalled || value != 123 {
			t.Errorf("Expected compute not called, value=123, got called=%v, value=%d", computeCalled, value)
		}
	})

	t.Run("update", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()
		cm.Set("counter", 10)

		// Update existing key
		updated := cm.Update("counter", func(oldValue int) int {
			return oldValue + 5
		})

		if !updated {
			t.Error("Expected update to succeed")
		}

		if value, _ := cm.Get("counter"); value != 15 {
			t.Errorf("Expected counter=15, got %d", value)
		}

		// Update non-existent key
		updated = cm.Update("nonexistent", func(oldValue int) int {
			return oldValue + 1
		})

		if updated {
			t.Error("Expected update of nonexistent key to fail")
		}
	})

	t.Run("compare and swap", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()
		cm.Set("value", 10)

		// Successful CAS
		swapped := cm.CompareAndSwap("value", 10, 20, func(current, expected int) bool {
			return current == expected
		})

		if !swapped {
			t.Error("Expected CAS to succeed")
		}

		if value, _ := cm.Get("value"); value != 20 {
			t.Errorf("Expected value=20, got %d", value)
		}

		// Failed CAS
		swapped = cm.CompareAndSwap("value", 10, 30, func(current, expected int) bool {
			return current == expected
		})

		if swapped {
			t.Error("Expected CAS to fail")
		}

		if value, _ := cm.Get("value"); value != 20 {
			t.Errorf("Expected value to remain 20, got %d", value)
		}
	})

	t.Run("iteration", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()
		expected := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}

		for k, v := range expected {
			cm.Set(k, v)
		}

		// Test ForEach
		found := make(map[string]int)
		cm.ForEach(func(key string, value int) bool {
			found[key] = value
			return true // continue
		})

		for k, v := range expected {
			if found[k] != v {
				t.Errorf("Expected %s=%d, got %d", k, v, found[k])
			}
		}

		// Test early termination
		count := 0
		cm.ForEach(func(key string, value int) bool {
			count++
			return count < 2 // stop after 2 items
		})

		if count != 2 {
			t.Errorf("Expected early termination at 2, got %d", count)
		}
	})

	t.Run("size and collections", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()
		
		if cm.Size() != 0 {
			t.Errorf("Expected size 0, got %d", cm.Size())
		}

		cm.Set("a", 1)
		cm.Set("b", 2)
		
		if cm.Size() != 2 {
			t.Errorf("Expected size 2, got %d", cm.Size())
		}

		keys := cm.Keys()
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}

		values := cm.Values()
		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}

		items := cm.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	})

	t.Run("filter", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()
		cm.Set("small", 5)
		cm.Set("medium", 15)
		cm.Set("large", 25)

		// Filter values > 10
		filtered := cm.Filter(func(key string, value int) bool {
			return value > 10
		})

		if filtered.Size() != 2 {
			t.Errorf("Expected 2 filtered items, got %d", filtered.Size())
		}

		if !filtered.Has("medium") || !filtered.Has("large") {
			t.Error("Expected medium and large to be in filtered map")
		}

		if filtered.Has("small") {
			t.Error("Expected small to not be in filtered map")
		}
	})

	t.Run("clear", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[string, int]()
		cm.Set("a", 1)
		cm.Set("b", 2)

		if cm.Size() != 2 {
			t.Error("Expected map to have items before clear")
		}

		cm.Clear()

		if cm.Size() != 0 {
			t.Error("Expected map to be empty after clear")
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		cm := patterns.NewConcurrentMap[int, int]()
		const numGoroutines = 100
		const opsPerGoroutine = 100

		var wg sync.WaitGroup
		
		// Multiple goroutines doing concurrent operations
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				for j := 0; j < opsPerGoroutine; j++ {
					key := id*opsPerGoroutine + j
					
					// Set
					cm.Set(key, key*2)
					
					// Get
					cm.Get(key)
					
					// Update
					cm.Update(key, func(oldValue int) int {
						return oldValue + 1
					})
					
					// GetOrSet
					cm.GetOrSet(key+10000, key)
				}
			}(i)
		}

		wg.Wait()

		// Verify some operations worked
		if cm.Size() == 0 {
			t.Error("Expected map to have items after concurrent operations")
		}
	})
}

func BenchmarkConcurrentMap(b *testing.B) {
	cm := patterns.NewConcurrentMap[int, int]()
	
	b.Run("Set", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				cm.Set(i, i*2)
				i++
			}
		})
	})
	
	b.Run("Get", func(b *testing.B) {
		// Pre-populate
		for i := 0; i < 1000; i++ {
			cm.Set(i, i*2)
		}
		
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				cm.Get(i % 1000)
				i++
			}
		})
	})
}