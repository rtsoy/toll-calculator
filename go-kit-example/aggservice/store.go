package aggservice

import (
	"fmt"
	"github.com/rtsoy/toll-calculator/types"
)

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() Storer {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (s *MemoryStore) Get(id int) (float64, error) {
	value, ok := s.data[id]
	if !ok {
		return 0, fmt.Errorf("could not find a distance for obuID=%d", id)
	}
	return value, nil
}

func (s *MemoryStore) Insert(distance types.Distance) error {
	s.data[distance.OBUID] += distance.Value
	return nil
}
