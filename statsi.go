// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package statsi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

var ErrLengthTooLong = errors.New("statsi: length too long to marshal")

const majorVersion = 0x01

type tag byte

const (
	tagCounter tag = iota
)

var now = time.Now // for testing

type Stats struct {
	mu       sync.Mutex
	last     time.Time
	counters []*Counter
}

func New() *Stats {
	return &Stats{
		last: now(),
	}
}

func (s *Stats) marshal() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var buf bytes.Buffer
	var varIntBuf [binary.MaxVarintLen64]byte

	now, last := now(), s.last
	s.last = now

	varIntLen := binary.PutUvarint(varIntBuf[:], uint64(now.Sub(last)/time.Millisecond))

	buf.Grow(1 + 2 + varIntLen)
	buf.WriteByte(majorVersion)
	binary.Write(&buf, binary.BigEndian, uint16(0)) // length placeholder
	buf.Write(varIntBuf[:varIntLen])

	headerLen := buf.Len()

	for _, c := range s.counters {
		varIntLen = binary.PutUvarint(varIntBuf[:], c.getDelta())

		length := varIntLen + len(c.wire)
		if length > 0x0fff {
			return nil, ErrLengthTooLong
		}

		buf.Grow(2 + length)
		binary.Write(&buf, binary.BigEndian, uint16(tagCounter)<<12|uint16(length))
		buf.Write(varIntBuf[:varIntLen])
		buf.WriteString(c.wire)
	}

	out := buf.Bytes()

	length := len(out) - headerLen
	if length > int(^uint16(0)) {
		return nil, ErrLengthTooLong
	}

	binary.BigEndian.PutUint16(out[1:], uint16(length))
	return out, nil
}
