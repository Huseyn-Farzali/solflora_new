package main

import "sync"

type SPState struct {
	mu    sync.RWMutex
	store map[Variable]float64
}

func NewSPState() *SPState {
	return &SPState{
		store: make(map[Variable]float64),
	}
}

func (s *SPState) Set(variable Variable, sp float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[variable] = sp
}

func (s *SPState) GetAll() map[Variable]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copyState := make(map[Variable]float64)
	for k, v := range s.store {
		copyState[k] = v
	}
	return copyState
}
