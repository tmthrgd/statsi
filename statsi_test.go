// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package statsi

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type mockTime struct {
	t int64
	d time.Duration
}

func (m *mockTime) Now() time.Time {
	return time.Unix(0, atomic.AddInt64(&m.t, int64(m.d)))
}

func TestMarshal(t *testing.T) {
	mock := &mockTime{0, 10 * time.Millisecond}
	now = mock.Now

	s := New()

	one := s.NewCounter("/test/counters/one")
	two := s.NewCounter("/test/counters/two")

	for i := 0; i < 3; i++ {
		s.NewCounter(fmt.Sprintf("/test/counters/n/%d", i))
	}

	for i := 1; i < 0x2000; i, mock.d = i<<1, mock.d*2 {
		one.Add(uint64(i))
		two.Add(uint64(i) / 2)

		b, err := s.marshal()
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%d:%x", len(b), b)
	}
}

func BenchmarkMarshal(b *testing.B) {
	now = (&mockTime{0, 10 * time.Millisecond}).Now

	s := New()
	s.NewCounter("/test/counters/one")

	for n := 0; n < b.N; n++ {
		if _, err := s.marshal(); err != nil {
			b.Fatal(err)
		}
	}
}
