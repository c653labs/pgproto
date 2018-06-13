package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type CloseCompleteTestSuite struct {
	suite.Suite
}

func TestCloseCompleteTestSuite(t *testing.T) {
	suite.Run(t, new(CloseCompleteTestSuite))
}

func (s *CloseCompleteTestSuite) Test_ParseCloseComplete() {
	raw := []byte{
		// Tag
		'3',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	close, err := pgproto.ParseCloseComplete(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(close)
	s.Equal(raw, close.Encode())
}

func BenchmarkCloseCompleteParse(b *testing.B) {
	raw := []byte{
		// Tag
		'3',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := pgproto.ParseCloseComplete(bytes.NewReader(raw))
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func (s *CloseCompleteTestSuite) Test_ParseCloseComplete_Empty() {
	close, err := pgproto.ParseCloseComplete(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(close)
}

func BenchmarkCloseCompleteParse_Empty(b *testing.B) {
	raw := []byte{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pgproto.ParseCloseComplete(bytes.NewReader(raw))
		}
	})
}

func (s *CloseCompleteTestSuite) Test_EncodeCloseComplete() {
	expected := []byte{
		// Tag
		'3',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	close := &pgproto.CloseComplete{}
	s.Equal(expected, close.Encode())
}

func BenchmarkCloseComplete(b *testing.B) {
	close := &pgproto.CloseComplete{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			close.Encode()
		}
	})
}
