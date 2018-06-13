package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type ParseCompleteTestSuite struct {
	suite.Suite
}

func TestParseCompleteTestSuite(t *testing.T) {
	suite.Run(t, new(ParseCompleteTestSuite))
}

func (s *ParseCompleteTestSuite) Test_ParseParseComplete() {
	raw := []byte{
		// Tag
		'1',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	parse, err := pgproto.ParseParseComplete(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(parse)
	s.Equal(raw, parse.Encode())
}

func BenchmarkParseCompleteParse(b *testing.B) {
	raw := []byte{
		// Tag
		'1',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := pgproto.ParseParseComplete(bytes.NewReader(raw))
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func (s *ParseCompleteTestSuite) Test_ParseParseComplete_Empty() {
	parse, err := pgproto.ParseParseComplete(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(parse)
}

func BenchmarkParseCompleteParse_Empty(b *testing.B) {
	raw := []byte{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pgproto.ParseParseComplete(bytes.NewReader(raw))
		}
	})
}

func (s *ParseCompleteTestSuite) Test_EncodeParseComplete() {
	expected := []byte{
		// Tag
		'1',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}

	parse := &pgproto.ParseComplete{}
	s.Equal(expected, parse.Encode())
}

func BenchmarkParseComplete(b *testing.B) {
	parse := &pgproto.ParseComplete{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			parse.Encode()
		}
	})
}
