package hashmap

import (
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
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

	var skeys []int
	s.Scan(func(key int) bool {
		skeys = append(skeys, key)
		return true
	})
	if len(skeys) != s.Len() {
		t.Fatalf("expected %d got %d", len(skeys), s.Len())
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

func copySetEntries(m *Set[int]) []int {
	all := m.Keys()
	sort.Ints(all)
	return all
}

func setEntriesEqual(a, b []int) bool {
	return reflect.DeepEqual(a, b)
}

func copySetTest(N int, s1 *Set[int], e11 []int, deep bool) {
	e12 := copySetEntries(s1)
	if !setEntriesEqual(e11, e12) {
		panic("!")
	}

	// Make a copy and compare the values
	s2 := s1.Copy()
	e21 := copySetEntries(s1)
	if !setEntriesEqual(e21, e12) {
		panic("!")
	}

	// Delete every other key
	var e22 []int
	for i, j := range rand.Perm(N) {
		if i&1 == 0 {

			e22 = append(e22, e21[j])
		} else {
			s2.Delete(e21[j])
		}
	}

	if s2.Len() != N/2 {
		panic("!")
	}
	sort.Ints(e22)
	e23 := copySetEntries(s2)
	if !setEntriesEqual(e23, e22) {
		panic("!")
	}
	if !deep {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			copySetTest(N/2, s2, e23, true)
		}()
		go func() {
			defer wg.Done()
			copySetTest(N/2, s2, e23, true)
		}()
		wg.Wait()
	}
	e24 := copySetEntries(s2)
	if !setEntriesEqual(e24, e23) {
		panic("!")
	}

}

func TestSetCopy(t *testing.T) {
	N := 1_000
	// create the initial map

	s1 := new(Set[int])
	for s1.Len() < N {
		s1.Insert(rand.Int())
	}
	e11 := copySetEntries(s1)
	dur := time.Second * 2
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			for time.Since(start) < dur {
				copySetTest(N, s1, e11, false)
			}
		}()
	}
	wg.Wait()
	e12 := copySetEntries(s1)
	if !setEntriesEqual(e11, e12) {
		panic("!")
	}
}

func TestSetGetPos(t *testing.T) {
	var m Set[int]
	if _, ok := m.GetPos(100); ok {
		t.Fatal()
	}
	for i := 0; i < 1000; i++ {
		m.Insert(i)
	}
	m2 := make(map[int]int)
	for i := 0; i < 10000; i++ {
		key, ok := m.GetPos(uint64(i))
		if !ok {
			t.Fatal()
		}
		m2[key] = key
	}
	if len(m2) != m.Len() {
		t.Fatal()
	}
}
