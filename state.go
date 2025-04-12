package main

import "sync"

type SPState struct {
	mu    sync.RWMutex
	store map[PhysicalVariable]float64
}

func NewSPState() *SPState {
	return &SPState{
		store: make(map[PhysicalVariable]float64),
	}
}

func (s *SPState) Set(variable PhysicalVariable, sp float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[variable] = sp
}

func (s *SPState) GetAll() map[PhysicalVariable]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copyState := make(map[PhysicalVariable]float64)
	for k, v := range s.store {
		copyState[k] = v
	}
	return copyState
}

type TuneState struct {
	mu    sync.RWMutex
	store map[PhysicalVariable]TuneProfile
}

func NewTuneState() *TuneState {
	return &TuneState{
		store: make(map[PhysicalVariable]TuneProfile),
	}
}

func (t *TuneState) Set(variable PhysicalVariable, value TuneProfile) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.store[variable] = value
}

func (t *TuneState) GetAll() map[PhysicalVariable]TuneProfile {
	t.mu.RLock()
	defer t.mu.RUnlock()

	copyState := make(map[PhysicalVariable]TuneProfile)
	for k, v := range t.store {
		copyState[k] = v
	}
	return copyState
}
