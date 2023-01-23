package suspender

import (
	"context"
	"errors"
	"fmt"
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

func New[T constraints.Ordered](cfg Config) (*Suspender[T], error) {
	if cfg.Count == 0 {
		cfg.Count = defaultCount
	}

	s := &Suspender[T]{
		m:   make(map[T]uint64),
		cfg: cfg,
	}

	return s, nil
}

func (s *Suspender[T]) IncWithCtx(ctx context.Context, value T) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("ctx.Err: %w", err)
	}
	if err := s.Inc(value); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		s.Dec(value)
	}()

	return nil
}

func (s *Suspender[T]) Inc(value T) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	v := s.m[value]
	if v >= s.cfg.Count {
		return ErrCountOverflow
	}

	s.m[value] = v + 1

	return nil
}

func (s *Suspender[T]) Dec(value T) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	count := s.m[value]
	if count > 0 {
		s.m[value] = count - 1
	}

	return nil
}
