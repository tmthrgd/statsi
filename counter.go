// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package statsi

import (
	"fmt"
	"sync/atomic"
)

type Counter struct {
	name, wire string
	val        uint64
}

func (s *Stats) NewCounter(name string) *Counter {
	s.mu.Lock()

	for _, c := range s.counters {
		if c.name == name {
			s.mu.Unlock()
			return c
		}
	}

	c := &Counter{name, compressName(name), 0}
	s.counters = append(s.counters, c)

	s.mu.Unlock()
	return c
}

func (s *Stats) GetCounter(name string) *Counter {
	s.mu.Lock()

	for _, c := range s.counters {
		if c.name == name {
			s.mu.Unlock()
			return c
		}
	}

	s.mu.Unlock()
	return nil
}

func (c *Counter) Increment() {
	atomic.AddUint64(&c.val, 1)
}

func (c *Counter) Add(incr uint64) {
	atomic.AddUint64(&c.val, incr)
}

func (c *Counter) getDelta() uint64 {
	return atomic.SwapUint64(&c.val, 0)
}

func (c *Counter) String() string {
	return fmt.Sprintf("*Counter{%s}", c.name)
}
