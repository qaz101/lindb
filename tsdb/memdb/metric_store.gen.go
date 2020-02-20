// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: int_map.tmpl

package memdb

import (
	"github.com/lindb/roaring"
)

// MetricStore represents int map using roaring bitmap
type MetricStore struct {
	keys   *roaring.Bitmap // store all keys
	values [][]tStoreINTF  // store all values by high/low key
}

// NewMetricStore creates a int map
func NewMetricStore() *MetricStore {
	return &MetricStore{
		keys: roaring.New(),
	}
}

// Get returns value by key, if exist returns it, else returns nil, false
func (m *MetricStore) Get(key uint32) (tStoreINTF, bool) {
	if len(m.values) == 0 {
		return nil, false
	}
	// get high index
	found, highIdx := m.keys.ContainsAndRankForHigh(key)
	if !found {
		return nil, false
	}
	// get low index
	found, lowIdx := m.keys.ContainsAndRankForLow(key, highIdx-1)
	if !found {
		return nil, false
	}
	return m.values[highIdx-1][lowIdx-1], true
}

// Put puts the value by key
func (m *MetricStore) Put(key uint32, value tStoreINTF) {
	if len(m.values) == 0 {
		// if values is empty, append new low container directly
		m.values = append(m.values, []tStoreINTF{value})

		m.keys.Add(key)
		return
	}
	found, highIdx := m.keys.ContainsAndRankForHigh(key)
	if !found {
		// high container not exist, insert it
		stores := m.values
		// insert operation, insert high values
		stores = append(stores, nil)
		copy(stores[highIdx+1:], stores[highIdx:len(stores)-1])
		stores[highIdx] = []tStoreINTF{value}
		m.values = stores

		m.keys.Add(key)
		return
	}
	// high container exist
	lowIdx := m.keys.RankForLow(key, highIdx-1)
	stores := m.values[highIdx-1]
	// insert operation
	stores = append(stores, nil)
	copy(stores[lowIdx+1:], stores[lowIdx:len(stores)-1])
	stores[lowIdx] = value
	m.values[highIdx-1] = stores

	m.keys.Add(key)
}

// Keys returns the all keys
func (m *MetricStore) Keys() *roaring.Bitmap {
	return m.keys
}

// Values returns the all values
func (m *MetricStore) Values() [][]tStoreINTF {
	return m.values
}

// size returns the size of keys
func (m *MetricStore) Size() int {
	return int(m.keys.GetCardinality())
}

// WalkEntry walks each kv entry via fn.
func (m *MetricStore) WalkEntry(fn func(key uint32, value tStoreINTF) error) error {
	values := m.values
	keys := m.keys
	highKeys := keys.GetHighKeys()
	for highIdx, highKey := range highKeys {
		hk := uint32(highKey) << 16
		lowValues := values[highIdx]
		lowContainer := keys.GetContainerAtIndex(highIdx)
		it := lowContainer.PeekableIterator()
		idx := 0
		for it.HasNext() {
			lowKey := it.Next()
			value := lowValues[idx]
			idx++
			if err := fn(uint32(lowKey&0xFFFF)|hk, value); err != nil {
				return err
			}
		}
	}
	return nil
}
