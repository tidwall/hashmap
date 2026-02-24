// Copyright 2019 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package hashmap

import (
	"unsafe"

	"github.com/zeebo/xxh3"
)

const (
	loadFactor  = 0.60                      // must be above 50%
	dibBitSize  = 16                        // 0xFFFF
	hashBitSize = 64 - dibBitSize           // 0xFFFFFFFFFFFF
	maxDIB      = ^uint64(0) >> hashBitSize // max 65,535
)

type entry[K comparable, V any] struct {
	hdib  uint64 // bitfield { hash:48 dib:16 }
	value V      // user value
	key   K      // user key
}

func (e *entry[K, V]) dib() int {
	return int(e.hdib & maxDIB)
}
func (e *entry[K, V]) hash() int {
	return int(e.hdib >> dibBitSize)
}
func (e *entry[K, V]) setDIB(dib int) {
	e.hdib = e.hdib>>dibBitSize<<dibBitSize | uint64(dib)&maxDIB
}
func (e *entry[K, V]) setHash(hash int) {
	e.hdib = uint64(hash)<<dibBitSize | e.hdib&maxDIB
}
func makeHDIB(hash, dib int) uint64 {
	return uint64(hash)<<dibBitSize | uint64(dib)&maxDIB
}

// https://zimbry.blogspot.com/2011/09/better-bit-mixing-improving-on.html
// hash u64 using mix13
func mix13(key uint64) uint64 {
	key ^= key >> 30
	key *= 0xbf58476d1ce4e5b9
	key ^= key >> 27
	key *= 0x94d049bb133111eb
	key ^= key >> 31
	return key
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 |
		~uint32 | ~uint64 | ~uintptr
}

func hashN[T integer](p unsafe.Pointer) int {
	return int(mix13(uint64(*(*T)(p))) >> dibBitSize)
}

func hashI[K comparable](k K) int   { return hashN[int](unsafe.Pointer(&k)) }
func hashI8[K comparable](k K) int  { return hashN[int8](unsafe.Pointer(&k)) }
func hashI16[K comparable](k K) int { return hashN[int16](unsafe.Pointer(&k)) }
func hashI32[K comparable](k K) int { return hashN[int32](unsafe.Pointer(&k)) }
func hashI64[K comparable](k K) int { return hashN[int64](unsafe.Pointer(&k)) }
func hashU[K comparable](k K) int   { return hashN[uint](unsafe.Pointer(&k)) }
func hashU8[K comparable](k K) int  { return hashN[uint8](unsafe.Pointer(&k)) }
func hashU16[K comparable](k K) int { return hashN[uint16](unsafe.Pointer(&k)) }
func hashU32[K comparable](k K) int { return hashN[uint32](unsafe.Pointer(&k)) }
func hashU64[K comparable](k K) int { return hashN[uint64](unsafe.Pointer(&k)) }
func hashS[K comparable](k K) int {
	return int(xxh3.HashString(*(*string)(unsafe.Pointer(&k))) >> dibBitSize)
}
func hashV[K comparable](k K) int {
	return hashS(unsafe.String((*byte)(unsafe.Pointer(&k)), unsafe.Sizeof(k)))
}

// Map is a hashmap. Like map[string]interface{}
type Map[K comparable, V any] struct {
	cap      int
	length   int
	mask     int
	growAt   int
	shrinkAt int
	buckets  []entry[K, V]
	hash     func(K) int
}

// New returns a new Map. Like map[string]interface{}
func New[K comparable, V any](cap int) *Map[K, V] {
	m := new(Map[K, V])
	m.cap = cap
	sz := 8
	for sz < m.cap {
		sz *= 2
	}
	if m.cap > 0 {
		m.cap = sz
	}
	m.buckets = make([]entry[K, V], sz)
	m.mask = len(m.buckets) - 1
	m.growAt = int(float64(len(m.buckets)) * loadFactor)
	m.shrinkAt = int(float64(len(m.buckets)) * (1 - loadFactor))

	var key K
	switch any(key).(type) {
	case int:
		m.hash = hashI
	case int8:
		m.hash = hashI8
	case int16:
		m.hash = hashI16
	case int32:
		m.hash = hashI32
	case int64:
		m.hash = hashI64
	case uint:
		m.hash = hashU
	case uint8:
		m.hash = hashU8
	case uint16:
		m.hash = hashU16
	case uint32:
		m.hash = hashU32
	case uint64:
		m.hash = hashU64
	case string:
		m.hash = hashS
	default:
		m.hash = hashV
	}
	return m
}

func (m *Map[K, V]) resize(newCap int) {
	nmap := New[K, V](newCap)
	for i := 0; i < len(m.buckets); i++ {
		if m.buckets[i].dib() > 0 {
			nmap.set(m.buckets[i].hash(), m.buckets[i].key, m.buckets[i].value)
		}
	}
	cap := m.cap
	*m = *nmap
	m.cap = cap
}

// Set assigns a value to a key.
// Returns the previous value, or false when no value was assigned.
func (m *Map[K, V]) Set(key K, value V) (V, bool) {
	if len(m.buckets) == 0 {
		*m = *New[K, V](0)
	}
	if m.length >= m.growAt {
		m.resize(len(m.buckets) * 2)
	}
	return m.set(m.hash(key), key, value)
}

func (m *Map[K, V]) set(hash int, key K, value V) (prev V, ok bool) {
	e := entry[K, V]{makeHDIB(hash, 1), value, key}
	i := e.hash() & m.mask
	for {
		b := &m.buckets[i]
		if b.dib() == 0 {
			*b = e
			m.length++
			return prev, false
		}
		if e.hash() == b.hash() && e.key == b.key {
			prev = b.value
			b.value = e.value
			return prev, true
		}
		if b.dib() < e.dib() {
			e, *b = *b, e
		}
		i = (i + 1) & m.mask
		e.setDIB(e.dib() + 1)
	}
}

// Get returns a value for a key.
// Returns false when no value has been assign for key.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	if len(m.buckets) == 0 {
		return value, false
	}
	hash := m.hash(key)
	i := hash & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			return value, false
		}
		if m.buckets[i].hash() == hash && m.buckets[i].key == key {
			return m.buckets[i].value, true
		}
		i = (i + 1) & m.mask
	}
}

// Len returns the number of values in map.
func (m *Map[K, V]) Len() int {
	return m.length
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *Map[K, V]) Delete(key K) (prev V, deleted bool) {
	if len(m.buckets) == 0 {
		return prev, false
	}
	hash := m.hash(key)
	i := hash & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			return prev, false
		}
		if m.buckets[i].hash() == hash && m.buckets[i].key == key {
			prev = m.buckets[i].value
			m.buckets[i].hdib = 0
			for {
				pi := i
				i = (i + 1) & m.mask
				if m.buckets[i].dib() <= 1 {
					m.buckets[pi] = entry[K, V]{}
					break
				}
				m.buckets[pi] = m.buckets[i]
				m.buckets[pi].setDIB(m.buckets[pi].dib() - 1)
			}
			m.length--
			if len(m.buckets) > m.cap && m.length <= m.shrinkAt {
				m.resize(m.length)
			}
			return prev, true
		}
		i = (i + 1) & m.mask
	}
}

// Scan iterates over all key/values.
// It's not safe to call or Set or Delete while scanning.
func (m *Map[K, V]) Scan(iter func(key K, value V) bool) {
	for i := 0; i < len(m.buckets); i++ {
		if m.buckets[i].dib() > 0 {
			if !iter(m.buckets[i].key, m.buckets[i].value) {
				return
			}
		}
	}
}

// Keys returns all keys as a slice
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, m.length)
	for i := 0; i < len(m.buckets); i++ {
		if m.buckets[i].dib() > 0 {
			keys = append(keys, m.buckets[i].key)
		}
	}
	return keys
}

// Values returns all values as a slice
func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, m.length)
	for i := 0; i < len(m.buckets); i++ {
		if m.buckets[i].dib() > 0 {
			values = append(values, m.buckets[i].value)
		}
	}
	return values
}

// Copy the hashmap.
func (m *Map[K, V]) Copy() *Map[K, V] {
	m2 := new(Map[K, V])
	*m2 = *m
	m2.buckets = make([]entry[K, V], len(m.buckets))
	copy(m2.buckets, m.buckets)
	return m2
}

// GetPos gets a single keys/value nearby a position.
// The pos param can be any valid uint64. Useful for grabbing a random item
// from the map.
func (m *Map[K, V]) GetPos(pos uint64) (key K, value V, ok bool) {
	for i := 0; i < len(m.buckets); i++ {
		index := (pos + uint64(i)) & uint64(m.mask)
		if m.buckets[index].dib() > 0 {
			return m.buckets[index].key, m.buckets[index].value, true
		}
	}
	// Empty map
	return key, value, false
}
