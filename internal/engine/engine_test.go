package engine

import (
	"fmt"
	"sync"
	"testing"
)

func TestEngine_SetGet(t *testing.T) {

	type tests struct {
		input []string
		want  string
	}

	tt := map[string]tests{
		"Simple Name Get and SET": {
			input: []string{"Name", "Drew"},
			want:  "Drew",
		},
	}
	kv := NewEngine()
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			kv.Set(tc.input[0], []byte(tc.input[1]))
			returned := string(kv.Get(tc.input[0]))
			if returned != tc.input[1] {
				t.Errorf("got = %v, want = %v", returned, tc.input[1])
			}
		})
	}
}

func TestEngine_Delete(t *testing.T) {

	type tests struct {
		key   string
		value string
		want  string
	}

	tt := map[string]tests{
		"Simple Delete key": {
			key:   "Name",
			value: "Drew",
			want:  "",
		},
	}
	kv := NewEngine()
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			kv.Set(tc.key, []byte(tc.value))
			kv.Del(tc.key)
			deleted := kv.Get(tc.key)
			if string(deleted) != tc.want {
				t.Errorf("It was not deleted correctly.")
			}

		})
	}

}

// Test for concurrency

func TestEngine_ConcurrentAccess(t *testing.T) {
	e := NewEngine()
	var wg sync.WaitGroup
	numGoroutines := 100

	wg.Add(numGoroutines * 2)

	// Writers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			e.Set(key, []byte("value"))
		}(i)
	}

	// Readers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			e.Get(key)
		}(i)
	}

	wg.Wait()
}
