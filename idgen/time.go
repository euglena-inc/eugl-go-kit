package idgen

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"time"
)

type Generator interface {
	NextID() int64
}

type TimeGenerator struct {
	clock  func() time.Time
	reader io.Reader
}

type TimeGeneratorOption func(*TimeGenerator)

func NewTimeGenerator(options ...TimeGeneratorOption) *TimeGenerator {
	generator := &TimeGenerator{
		clock:  time.Now,
		reader: rand.Reader,
	}
	for _, option := range options {
		option(generator)
	}
	if generator.clock == nil {
		generator.clock = time.Now
	}
	if generator.reader == nil {
		generator.reader = rand.Reader
	}
	return generator
}

func WithClock(clock func() time.Time) TimeGeneratorOption {
	return func(generator *TimeGenerator) {
		generator.clock = clock
	}
}

func WithReader(reader io.Reader) TimeGeneratorOption {
	return func(generator *TimeGenerator) {
		generator.reader = reader
	}
}

func (g *TimeGenerator) NextID() int64 {
	var random [2]byte
	if _, err := io.ReadFull(g.reader, random[:]); err != nil {
		return g.clock().UnixMilli() * 10000
	}
	suffix := int64(binary.BigEndian.Uint16(random[:]) % 10000)
	return g.clock().UnixMilli()*10000 + suffix
}
