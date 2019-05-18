// Copyright 2019 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package rhh

import (
	"reflect"
	"unsafe"

	"github.com/cespare/xxhash"
)

type entryU64 struct {
	hdib  uint64      // bitfield { hash:48 dib:16 }
	key   uint64      // user key
	value interface{} // user value
}

func (e *entryU64) dib() int {
	return int(e.hdib & maxDIB)
}
func (e *entryU64) hash() int {
	return int(e.hdib >> dibBitSize)
}
func (e *entryU64) setDIB(dib int) {
	e.hdib = e.hdib>>dibBitSize<<dibBitSize | uint64(dib)&maxDIB
}
func (e *entryU64) setHash(hash int) {
	e.hdib = uint64(hash)<<dibBitSize | e.hdib&maxDIB
}

// hash returns a 48-bit hash for 64-bit environments, or 32-bit hash for
// 32-bit environments.
func (m *MapU64) hash(key uint64) int {
	return int(xxhash.Sum64(*(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&key)), Len: 8, Cap: 8,
	}))) >> dibBitSize)
}

// MapU64 is a map. Like map[uint64]interface{}
type MapU64 struct {
	cap      int
	length   int
	mask     int
	growAt   int
	shrinkAt int
	buckets  []entryU64
}

// NewU64 returns a new map. Like map[uint64]interface{}
func NewU64(cap int) *MapU64 {
	m := new(MapU64)
	m.cap = cap
	sz := 8
	for sz < m.cap {
		sz *= 2
	}
	m.buckets = make([]entryU64, sz)
	m.mask = len(m.buckets) - 1
	m.growAt = int(float64(len(m.buckets)) * loadFactor)
	m.shrinkAt = int(float64(len(m.buckets)) * (1 - loadFactor))
	return m
}

func (m *MapU64) resize(newCap int) {
	nmap := NewU64(newCap)
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
func (m *MapU64) Set(key uint64, value interface{}) (interface{}, bool) {
	if len(m.buckets) == 0 {
		*m = *NewU64(0)
	}
	if m.length >= m.growAt {
		m.resize(len(m.buckets) * 2)
	}
	return m.set(m.hash(key), key, value)
}

func (m *MapU64) set(hash int, key uint64, value interface{}) (interface{}, bool) {
	e := entryU64{makeHDIB(hash, 1), key, value}
	i := e.hash() & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			m.buckets[i] = e
			m.length++
			return nil, false
		}
		if e.hash() == m.buckets[i].hash() && e.key == m.buckets[i].key {
			old := m.buckets[i].value
			m.buckets[i].value = e.value
			return old, true
		}
		if m.buckets[i].dib() < e.dib() {
			e, m.buckets[i] = m.buckets[i], e
		}
		i = (i + 1) & m.mask
		e.setDIB(e.dib() + 1)
	}
}

// Get returns a value for a key.
// Returns false when no value has been assign for key.
func (m *MapU64) Get(key uint64) (interface{}, bool) {
	if len(m.buckets) == 0 {
		return nil, false
	}
	hash := m.hash(key)
	i := hash & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			return nil, false
		}
		if m.buckets[i].hash() == hash && m.buckets[i].key == key {
			return m.buckets[i].value, true
		}
		i = (i + 1) & m.mask
	}
}

// Len returns the number of values in map.
func (m *MapU64) Len() int {
	return m.length
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *MapU64) Delete(key uint64) (interface{}, bool) {
	if len(m.buckets) == 0 {
		return nil, false
	}
	hash := m.hash(key)
	i := hash & m.mask
	for {
		if m.buckets[i].dib() == 0 {
			return nil, false
		}
		if m.buckets[i].hash() == hash && m.buckets[i].key == key {
			old := m.buckets[i].value
			m.remove(i)
			return old, true
		}
		i = (i + 1) & m.mask
	}
}

func (m *MapU64) remove(i int) {
	m.buckets[i].setDIB(0)
	for {
		pi := i
		i = (i + 1) & m.mask
		if m.buckets[i].dib() <= 1 {
			m.buckets[pi] = entryU64{}
			break
		}
		m.buckets[pi] = m.buckets[i]
		m.buckets[pi].setDIB(m.buckets[pi].dib() - 1)
	}
	m.length--
	if len(m.buckets) > m.cap && m.length <= m.shrinkAt {
		m.resize(m.length)
	}
}

// Range iterates overall all key/values.
// It's not safe to call or Set or Delete while ranging.
func (m *MapU64) Range(iter func(key uint64, value interface{}) bool) {
	for i := 0; i < len(m.buckets); i++ {
		if m.buckets[i].dib() > 0 {
			if !iter(m.buckets[i].key, m.buckets[i].value) {
				return
			}
		}
	}
}
