package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type EmptyQueryResponseTestSuite struct {
	suite.Suite
}

func TestEmptyQueryResponseTestSuite(t *testing.T) {
	suite.Run(t, new(EmptyQueryResponseTestSuite))
}

func (s *EmptyQueryResponseTestSuite) Test_ParseEmptyQueryResponse() {
	raw := []byte{
		// Tag
		'I',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	empty, err := pgproto.ParseEmptyQueryResponse(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(empty)
	s.Equal(raw, empty.Encode())
}

func BenchmarkEmptyQueryResponseParse(b *testing.B) {
	raw := []byte{
		// Tag
		'I',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := pgproto.ParseEmptyQueryResponse(bytes.NewReader(raw))
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func (s *EmptyQueryResponseTestSuite) Test_ParseEmptyQueryResponse_Empty() {
	empty, err := pgproto.ParseEmptyQueryResponse(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(empty)
}

func BenchmarkEmptyQueryResponseParse_Empty(b *testing.B) {
	raw := []byte{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pgproto.ParseEmptyQueryResponse(bytes.NewReader(raw))
		}
	})
}

func (s *EmptyQueryResponseTestSuite) Test_EncodeEmptyQueryResponse() {
	expected := []byte{
		// Tag
		'I',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	empty := &pgproto.EmptyQueryResponse{}
	s.Equal(expected, empty.Encode())
}

func BenchmarkEmptyQueryResponseEncode(b *testing.B) {
	empty := &pgproto.EmptyQueryResponse{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			empty.Encode()
		}
	})
}
