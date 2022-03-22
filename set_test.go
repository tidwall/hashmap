package hashmap

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	var s Set[int]

	keys := rand.Perm(1000000)

	for i := 0; i < len(keys); i++ {
		s.Insert(keys[i])
		if s.Len() != i+1 {
			t.Fatalf("expected %d got %d", i+1, s.Len())
		}
	}

	for i := 0; i < len(keys); i++ {
		ok := s.Contains(keys[i])
		if !ok {
			t.Fatalf("expected true")
		}
	}

	for i := 0; i < len(keys); i++ {
		s.Delete(keys[i])
		if s.Len() != len(keys)-i-1 {
			t.Fatalf("expected %d got %d", len(keys)-i-1, s.Len())
		}
	}
}

func TestSetKeys(t *testing.T) {
	var s Set[string]
	s.Insert("key")
	expect := []string{"key"}
	got := s.Keys()
	if !reflect.DeepEqual(got, expect) {
		t.Fatal("expected Keys equal")
	}
}
