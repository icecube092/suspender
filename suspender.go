package suspender

import (
	"errors"
	"sync"

	"golang.org/x/exp/constraints"
)

var ErrCountOverflow = errors.New("count overflow")

type Suspender[T constraints.Ordered] struct {
	m   map[T]uint64
	mux sync.RWMutex

	cfg Config
}

const defaultCount uint64 = 1

func New[T constraints.Ordered](cfg Config) *Suspender[T] {
	if cfg.Count == 0 {
		cfg.Count = defaultCount
	}

	s := &Suspender[T]{
		m:   make(map[T]uint64),
		cfg: cfg,
	}

	return s
}

func (s *Suspender[T]) Add(value T) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	v := s.m[value]
	if v >= s.cfg.Count {
		return ErrCountOverflow
	}

	s.m[value] = v + 1

	return nil
}

func (s *Suspender[T]) Drop(value T) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.m[value] = 0

	return nil
}
