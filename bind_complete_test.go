package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type BindCompleteTestSuite struct {
	suite.Suite
}

func TestBindCompleteTestSuite(t *testing.T) {
	suite.Run(t, new(BindCompleteTestSuite))
}

func (s *BindCompleteTestSuite) Test_ParseBindComplete() {
	raw := []byte{
		// Tag
		'2',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	bind, err := pgproto.ParseBindComplete(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(bind)
	s.Equal(raw, bind.Encode())
}

func BenchmarkBindCompleteParse(b *testing.B) {
	raw := []byte{
		// Tag
		'2',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := pgproto.ParseBindComplete(bytes.NewReader(raw))
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func (s *BindCompleteTestSuite) Test_ParseBindComplete_Empty() {
	bind, err := pgproto.ParseBindComplete(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(bind)
}

func BenchmarkBindCompleteParse_Empty(b *testing.B) {
	raw := []byte{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pgproto.ParseBindComplete(bytes.NewReader(raw))
		}
	})
}

func (s *BindCompleteTestSuite) Test_EncodeBindComplete() {
	expected := []byte{
		// Tag
		'2',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	bind := &pgproto.BindComplete{}
	s.Equal(expected, bind.Encode())
}

func BenchmarkBindComplete(b *testing.B) {
	bind := &pgproto.BindComplete{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			bind.Encode()
		}
	})
}
