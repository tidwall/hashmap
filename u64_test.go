// Copyright 2019 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package rhh

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

type keyTU64 = uint64
type valueTU64 = interface{}

func kU64(key int) keyTU64 {
	return uint64(key)
}

func addU64(x keyTU64, delta int) int {
	return int(x) + delta
}

///////////////////////////
func randomU64(N int, perm bool) []keyTU64 {
	nums := make([]keyTU64, N)
	if perm {
		for i, x := range rand.Perm(N) {
			nums[i] = kU64(x)
		}
	} else {
		m := make(map[keyTU64]bool)
		for len(m) < N {
			m[kU64(int(rand.Uint64()))] = true
		}
		var i int
		for k := range m {
			nums[i] = k
			i++
		}
	}
	return nums
}

func shuffleU64(nums []keyTU64) {
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

func TestRandomDataU64(t *testing.T) {
	N := 10000
	start := time.Now()
	for time.Since(start) < time.Second*2 {
		nums := randomU64(N, true)
		var m *MapU64
		switch rand.Int() % 5 {
		default:
			m = NewU64(N / ((rand.Int() % 3) + 1))
		case 1:
			m = new(MapU64)
		case 2:
			m = NewU64(0)
		}
		v, ok := m.Get(kU64(999))
		if ok || v != nil {
			t.Fatalf("expected %v, got %v", nil, v)
		}
		v, ok = m.Delete(kU64(999))
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
		shuffleU64(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			if !ok || v == nil || v != nums[i] {
				t.Fatalf("expected %v, got %v", nums[i], v)
			}
		}
		// replace all the items
		shuffleU64(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Set(nums[i], addU64(nums[i], 1))
			if !ok || v != nums[i] {
				t.Fatalf("expected %v, got %v", nums[i], v)
			}
		}
		if m.Len() != N {
			t.Fatalf("expected %v, got %v", N, m.Len())
		}
		// retrieve all the items
		shuffleU64(nums)
		for i := 0; i < len(nums); i++ {
			v, ok := m.Get(nums[i])
			if !ok || v != addU64(nums[i], 1) {
				t.Fatalf("expected %v, got %v", addU64(nums[i], 1), v)
			}
		}
		// remove half the items
		shuffleU64(nums)
		for i := 0; i < len(nums)/2; i++ {
			v, ok := m.Delete(nums[i])
			if !ok || v != addU64(nums[i], 1) {
				t.Fatalf("expected %v, got %v", addU64(nums[i], 1), v)
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
			if !ok || v != addU64(nums[i], 1) {
				t.Fatalf("expected %v, got %v", addU64(nums[i], 1), v)
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
		m.Range(func(key keyTU64, value valueTU64) bool {
			if value != addU64(key, 1) {
				t.Fatalf("expected %v, got %v", addU64(key, 1), value)
			}
			return true
		})
		var n int
		m.Range(func(key keyTU64, value valueTU64) bool {
			n++
			return false
		})
		if n != 1 {
			t.Fatalf("expected %v, got %v", 1, n)
		}
		for i := len(nums) / 2; i < len(nums); i++ {
			v, ok := m.Delete(nums[i])
			if !ok || v != addU64(nums[i], 1) {
				t.Fatalf("expected %v, got %v", addU64(nums[i], 1), v)
			}
		}
	}
}

func TestBenchU64(t *testing.T) {
	N, _ := strconv.ParseUint(os.Getenv("MAPBENCH"), 10, 64)
	if N == 0 {
		fmt.Printf("Enable benchmarks with MAPBENCH=1000000\n")
		return
	}
	nums := randomU64(int(N), false)
	var pnums []valueTU64
	for i := range nums {
		pnums = append(pnums, valueTU64(&nums[i]))
	}
	fmt.Printf("\n## UINT64 KEYS\n\n")
	t.Run("RobinHood", func(t *testing.T) {
		testPerfU64(nums, pnums, "robinhood")
	})
	t.Run("Stdlib", func(t *testing.T) {
		testPerfU64(nums, pnums, "stdlib")
	})
}

func testPerfU64(nums []keyTU64, pnums []valueTU64, which string) {
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
		m := make(map[keyTU64]valueTU64, initSize)
		setop = func(i, _ int) { m[nums[i]] = pnums[i] }
		getop = func(i, _ int) { _ = m[nums[i]] }
		delop = func(i, _ int) { delete(m, nums[i]) }
		scnop = func() {
			for range m {
			}
		}
	case "robinhood":
		m := NewU64(initSize)
		setop = func(i, _ int) { m.Set(nums[i], pnums[i]) }
		getop = func(i, _ int) { m.Get(nums[i]) }
		delop = func(i, _ int) { m.Delete(nums[i]) }
		scnop = func() {
			m.Range(func(key keyTU64, value valueTU64) bool {
				return true
			})
		}
	}
	fmt.Printf("-- %s --", which)
	fmt.Printf("\n")

	ops := []func(int, int){setop, getop, setop, nil, delop}
	tags := []string{"set", "get", "reset", "scan", "delete"}
	for i := range ops {
		shuffleU64(nums)
		var na bool
		var n int
		start := time.Now()
		if tags[i] == "scan" {
			op := scnop
			if op == nil {
				na = true
			} else {
				n = 20
				lotsa.Ops(n, 1, func(_, _ int) { op() })
			}

		} else {
			n = len(nums)
			lotsa.Ops(n, 1, ops[i])
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

func TestHashDIBU64(t *testing.T) {
	var e entryU64
	e.setDIB(100)
	e.setHash(90000)
	if e.dib() != 100 {
		t.Fatalf("expected %v, got %v", 100, e.dib())
	}
	if e.hash() != 90000 {
		t.Fatalf("expected %v, got %v", 90000, e.hash())
	}
}
