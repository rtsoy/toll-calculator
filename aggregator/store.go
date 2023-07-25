package main

import "github.com/rtsoy/toll-calculator/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Insert(distance types.Distance) error {
	return nil
}
