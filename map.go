// Copyright 2019 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package rhh

import (
	"github.com/cespare/xxhash"
)

const (
	loadFactor  = 0.85                      // must be above 50%
	dibBitSize  = 16                        // 0xFFFF
	hashBitSize = 64 - dibBitSize           // 0xFFFFFFFFFFFF
	maxHash     = ^uint64(0) >> dibBitSize  // max 28,147,497,671,0655
	maxDIB      = ^uint64(0) >> hashBitSize // max 65,535
)

type entry struct {
	hdib  uint64      // bitfield { hash:48 dib:16 }
	key   string      // user key
	value interface{} // user value
}

func (e *entry) dib() int {
	return int(e.hdib & maxDIB)
}
func (e *entry) hash() int {
	return int(e.hdib >> dibBitSize)
}
func (e *entry) setDIB(dib int) {
	e.hdib = e.hdib>>dibBitSize<<dibBitSize | uint64(dib)&maxDIB
}
func (e *entry) setHash(hash int) {
	e.hdib = uint64(hash)<<dibBitSize | e.hdib&maxDIB
}
func makeHDIB(hash, dib int) uint64 {
	return uint64(hash)<<dibBitSize | uint64(dib)&maxDIB
}

// hash returns a 48-bit hash for 64-bit environments, or 32-bit hash for
// 32-bit environments.
func (m *Map) hash(key string) int {
	return int(xxhash.Sum64String(key) >> dibBitSize)
}

// Map is a hashmap. Like map[string]interface{}
type Map struct {
	cap      int
	length   int
	mask     int
	growAt   int
	shrinkAt int
	buckets  []entry
}

// New returns a new Map. Like map[string]interface{}
func New(cap int) *Map {
	m := new(Map)
	m.cap = cap
	sz := 8
	for sz < m.cap {
		sz *= 2
	}
	m.buckets = make([]entry, sz)
	m.mask = len(m.buckets) - 1
	m.growAt = int(float64(len(m.buckets)) * loadFactor)
	m.shrinkAt = int(float64(len(m.buckets)) * (1 - loadFactor))
	return m
}

func (m *Map) resize(newCap int) {
	nmap := New(newCap)
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
func (m *Map) Set(key string, value interface{}) (interface{}, bool) {
	if len(m.buckets) == 0 {
		*m = *New(0)
	}
	if m.length >= m.growAt {
		m.resize(len(m.buckets) * 2)
	}
	return m.set(m.hash(key), key, value)
}

func (m *Map) set(hash int, key string, value interface{}) (interface{}, bool) {
	e := entry{makeHDIB(hash, 1), key, value}
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
func (m *Map) Get(key string) (interface{}, bool) {
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
func (m *Map) Len() int {
	return m.length
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *Map) Delete(key string) (interface{}, bool) {
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

func (m *Map) remove(i int) {
	m.buckets[i].setDIB(0)
	for {
		pi := i
		i = (i + 1) & m.mask
		if m.buckets[i].dib() <= 1 {
			m.buckets[pi] = entry{}
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

// Range iterates over all key/values.
// It's not safe to call or Set or Delete while ranging.
func (m *Map) Range(iter func(key string, value interface{}) bool) {
	for i := 0; i < len(m.buckets); i++ {
		if m.buckets[i].dib() > 0 {
			if !iter(m.buckets[i].key, m.buckets[i].value) {
				return
			}
		}
	}
}

// GetPos gets a single keys/value nearby a position
// The pos param can be any valid uint64. Useful for grabbing a random item
// from the map.
// It's not safe to call or Set or Delete while ranging.
func (m *Map) GetPos(pos uint64) (key string, value interface{}, ok bool) {
	for i := 0; i < len(m.buckets); i++ {
		index := (pos + uint64(i)) & uint64(m.mask)
		if m.buckets[index].dib() > 0 {
			return m.buckets[index].key, m.buckets[index].value, true
		}
	}
	return "", nil, false
}
