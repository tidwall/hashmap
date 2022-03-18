// Copyright 2019 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package hashmap

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type keyT = string
type valueT = interface{}

func k(key int) keyT {
	return strconv.FormatInt(int64(key), 10)
}

func add(x keyT, delta int) int {
	i, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(i + int64(delta))
}

///////////////////////////
func random(N int, perm bool) []keyT {
	nums := make([]keyT, N)
	if perm {
		for i, x := range rand.Perm(N) {
			nums[i] = k(x)
		}
	} else {
		m := make(map[keyT]bool)
		for len(m) < N {
			m[k(int(rand.Uint64()))] = true
		}
		var i int
		for k := range m {
			nums[i] = k
			i++
		}
	}
	return nums
}

func shuffle[K comparable](nums []K) {
	for i := range nums {
		j := rand.Intn(i + 1)
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func init() {
	//var seed int64 = 1519776033517775607
	seed := (time.Now().UnixNano())
	println("seed:", seed)
	rand.Seed(seed)
}

type imap struct {
	m *Map[string, interface{}]
}

func newimap(cap int) *imap {
	m := new(imap)
	m.m = New[string, interface{}](cap)
	return m
}

func (m *imap) Get(key string) (interface{}, bool) {
	if m.m == nil {
		return nil, false
	}
	return m.m.Get(key)
}
func (m *imap) Set(key string, value interface{}) (interface{}, bool) {
	if m.m == nil {
		m.m = new(Map[string, interface{}])
	}
	return m.m.Set(key, value)
}
func (m *imap) Delete(key string) (interface{}, bool) {
	if m.m == nil {
		return nil, false
	}
	return m.m.Delete(key)
}
func (m *imap) Len() int {
	if m.m == nil {
		return 0
	}
	return m.m.Len()
}
func (m *imap) Scan(iter func(key string, value interface{}) bool) {
	if m.m == nil {
		return
	}
	m.m.Scan(iter)
}

func TestRandomData(t *testing.T) {
	N := 10000
	start := time.Now()
	for time.Since(start) < time.Second*2 {
		nums := random(N, true)
		var m *imap
		switch rand.Int() % 5 {
		default:
			m = newimap(N / ((rand.Int() % 3) + 1))
		case 1:
			m = new(imap)
		case 2:
			m = newimap(0)
		}
		v, ok := m.Get(k(999))
		if ok || v != nil {
			t.Fatalf("expected %v, got %v", nil, v)
		}
		v, ok = m.Delete(k(999))
		if ok || v != nil {
			t.Fatalf("expected %v, got %v", nil, v)
		}
		if m.Len() != 0 {
			t.Fatalf("expected %v, got %v", 0, m.Len())
		}
		// set a bunch of items
		for i := 0; i < len(nums); i++ {
			v, ok := m.Set(nums[i], nums[i])
			if ok || v != nil {
				t.Fatalf("expected %v, got %v", nil, v)
			}
		}
		if m.Len() != N {
			t.Fatalf("expected %v, got %v", N, m.Len())
		}
		// retrieve all the items
		shuffle(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			if !ok || v == nil || v != nums[i] {
				t.Fatalf("expected %v, got %v", nums[i], v)
			}
		}
		// replace all the items
		shuffle(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Set(nums[i], add(nums[i], 1))
			if !ok || v != nums[i] {
				t.Fatalf("expected %v, got %v", nums[i], v)
			}
		}
		if m.Len() != N {
			t.Fatalf("expected %v, got %v", N, m.Len())
		}
		// retrieve all the items
		shuffle(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
		// remove half the items
		shuffle(nums)
		for i := 0; i < len(nums)/2; i++ {
			v, ok := m.Delete(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
		if m.Len() != N/2 {
			t.Fatalf("expected %v, got %v", N/2, m.Len())
		}
		// check to make sure that the items have been removed
		for i := 0; i < len(nums)/2; i++ {
			v, ok := m.Get(nums[i])
			if ok || v != nil {
				t.Fatalf("expected %v, got %v", nil, v)
			}
		}
		// check the second half of the items
		for i := len(nums) / 2; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
		// try to delete again, make sure they don't exist
		for i := 0; i < len(nums)/2; i++ {
			v, ok := m.Delete(nums[i])
			if ok || v != nil {
				t.Fatalf("expected %v, got %v", nil, v)
			}
		}
		if m.Len() != N/2 {
			t.Fatalf("expected %v, got %v", N/2, m.Len())
		}
		m.Scan(func(key keyT, value valueT) bool {
			if value != add(key, 1) {
				t.Fatalf("expected %v, got %v", add(key, 1), value)
			}
			return true
		})
		var n int
		m.Scan(func(key keyT, value valueT) bool {
			n++
			return false
		})
		if n != 1 {
			t.Fatalf("expected %v, got %v", 1, n)
		}
		for i := len(nums) / 2; i < len(nums); i++ {
			v, ok := m.Delete(nums[i])
			if !ok || v != add(nums[i], 1) {
				t.Fatalf("expected %v, got %v", add(nums[i], 1), v)
			}
		}
	}
}

func TestBench(t *testing.T) {
	N, _ := strconv.ParseUint(os.Getenv("MAPBENCH"), 10, 64)
	if N == 0 {
		fmt.Printf("Enable benchmarks with MAPBENCH=1000000\n")
		return
	}

	var pnums []int
	for i := 0; i < int(N); i++ {
		pnums = append(pnums, i)
	}

	{
		fmt.Printf("\n## STRING KEYS\n\n")
		nums := random(int(N), false)
		t.Run("Tidwall", func(t *testing.T) {
			testPerf(nums, pnums, "tidwall")
		})
		t.Run("Stdlib", func(t *testing.T) {
			testPerf(nums, pnums, "stdlib")
		})
	}
	{
		fmt.Printf("\n## INT KEYS\n\n")
		nums := rand.Perm(int(N))
		t.Run("Tidwall", func(t *testing.T) {
			testPerf(nums, pnums, "tidwall")
		})
		t.Run("Stdlib", func(t *testing.T) {
			testPerf(nums, pnums, "stdlib")
		})
	}

}

func printItem(s string, size int, dir int) {
	for len(s) < size {
		if dir == -1 {
			s += " "
		} else {
			s = " " + s
		}
	}
	fmt.Printf("%s ", s)
}

func testPerf[K comparable, V any](nums []K, pnums []V, which string) {
	var ms1, ms2 runtime.MemStats
	initSize := 0 //len(nums) * 2
	defer func() {
		heapBytes := int(ms2.HeapAlloc - ms1.HeapAlloc)
		fmt.Printf("memory %13s bytes %19s/entry \n",
			commaize(heapBytes), commaize(heapBytes/len(nums)))
		fmt.Printf("\n")
	}()
	runtime.GC()
	time.Sleep(time.Millisecond * 100)
	runtime.ReadMemStats(&ms1)

	var setop, getop, delop func(int, int)
	var scnop func()
	switch which {
	case "stdlib":
		m := make(map[K]V, initSize)
		setop = func(i, _ int) { m[nums[i]] = pnums[i] }
		getop = func(i, _ int) { _ = m[nums[i]] }
		delop = func(i, _ int) { delete(m, nums[i]) }
		scnop = func() {
			for range m {
			}
		}
	case "tidwall":
		var m Map[K, V]
		setop = func(i, _ int) { m.Set(nums[i], pnums[i]) }
		getop = func(i, _ int) { m.Get(nums[i]) }
		delop = func(i, _ int) { m.Delete(nums[i]) }
		scnop = func() {
			m.Scan(func(key K, value V) bool {
				return true
			})
		}
	}
	fmt.Printf("-- %s --", which)
	fmt.Printf("\n")

	ops := []func(int, int){setop, getop, setop, nil, delop}
	tags := []string{"set", "get", "reset", "scan", "delete"}
	for i := range ops {
		shuffle(nums)
		var na bool
		var n int
		start := time.Now()
		if tags[i] == "scan" {
			op := scnop
			if op == nil {
				na = true
			} else {
				n = 20
				for i := 0; i < n; i++ {
					op()
				}
			}
		} else {
			n = len(nums)
			for j := 0; j < n; j++ {
				ops[i](j, 1)
			}
		}
		dur := time.Since(start)
		if i == 0 {
			runtime.GC()
			time.Sleep(time.Millisecond * 100)
			runtime.ReadMemStats(&ms2)
		}
		printItem(tags[i], 9, -1)
		if na {
			printItem("-- unavailable --", 14, 1)
		} else {
			if n == -1 {
				printItem("unknown ops", 14, 1)
			} else {
				printItem(fmt.Sprintf("%s ops", commaize(n)), 14, 1)
			}
			printItem(fmt.Sprintf("%.0fms", dur.Seconds()*1000), 8, 1)
			if n != -1 {
				printItem(fmt.Sprintf("%s/sec", commaize(int(float64(n)/dur.Seconds()))), 18, 1)
			}
		}
		fmt.Printf("\n")
	}
}

func commaize(n int) string {
	s1, s2 := fmt.Sprintf("%d", n), ""
	for i, j := len(s1)-1, 0; i >= 0; i, j = i-1, j+1 {
		if j%3 == 0 && j != 0 {
			s2 = "," + s2
		}
		s2 = string(s1[i]) + s2
	}
	return s2
}

func TestHashDIB(t *testing.T) {
	var e entry[string, interface{}]
	e.setDIB(100)
	e.setHash(90000)
	if e.dib() != 100 {
		t.Fatalf("expected %v, got %v", 100, e.dib())
	}
	if e.hash() != 90000 {
		t.Fatalf("expected %v, got %v", 90000, e.hash())
	}
}

func TestIntInt(t *testing.T) {
	var m Map[int, int]

	keys := rand.Perm(1000000)

	for i := 0; i < len(keys); i++ {
		_, ok := m.Set(keys[i], keys[i]*10)
		if ok {
			t.Fatalf("expected false")
		}
		if m.Len() != i+1 {
			t.Fatalf("expected %d got %d", i+1, m.Len())
		}
	}

	for i := 0; i < len(keys); i++ {
		v, ok := m.Get(keys[i])
		if !ok {
			t.Fatalf("expected true")
		}
		if v != keys[i]*10 {
			t.Fatalf("expected %d got %d", keys[i]*10, v)
		}
	}

	for i := 0; i < len(keys); i++ {
		v, ok := m.Delete(keys[i])
		if !ok {
			t.Fatalf("expected true")
		}
		if v != keys[i]*10 {
			t.Fatalf("expected %d got %d", keys[i]*10, v)
		}
		if m.Len() != len(keys)-i-1 {
			t.Fatalf("expected %d got %d", len(keys)-i-1, m.Len())
		}
	}
}

func TestMapValues(t *testing.T) {
	var m Map[int, int]
	m.Set(1, 2)
	expect := []int{2}
	got := m.Values()
	if !reflect.DeepEqual(got, expect) {
		t.Fatal("expected Values equal")
	}
}
