package main

import (
	"fmt"

	"github.com/souvik150/supplychain-hackathon-proj/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	dist, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("no distance found for OBUID %d", id)
	}
	return dist, nil
}