package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

// 'X' [int32 - length]
var rawTerminationMessage = [5]byte{
	// Tag
	'X',
	// Length
	'\x00', '\x00', '\x00', '\x04',
}

type TerminationTestSuite struct {
	suite.Suite
}

func TestTerminationTestSuite(t *testing.T) {
	suite.Run(t, new(TerminationTestSuite))
}

func (s *TerminationTestSuite) Test_ParseTermination() {
	term, err := pgproto.ParseTermination(bytes.NewReader(rawTerminationMessage[:]))
	s.Nil(err)
	s.NotNil(term)
	s.Equal(rawTerminationMessage[:], term.Encode())
}

func BenchmarkParseTermination(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := pgproto.ParseTermination(bytes.NewReader(rawTerminationMessage[:]))
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func (s *TerminationTestSuite) Test_ParseTermination_Empty() {
	term, err := pgproto.ParseTermination(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(term)
}

func (s *TerminationTestSuite) Test_ParseTermination_InvalidTag() {
	// lowercase 'x' instead of uppercase 'X'
	raw := []byte{'x', '\x00', '\x00', '\x00', '\x04'}

	term, err := pgproto.ParseTermination(bytes.NewReader(raw))
	s.NotNil(err)
	s.Nil(term)
}

func (s *TerminationTestSuite) Test_ParseTermination_InvalidLength() {
	// length of 3 instead of 4
	raw := []byte{'X', '\x00', '\x00', '\x00', '\x03'}

	term, err := pgproto.ParseTermination(bytes.NewReader(raw))
	s.NotNil(err)
	s.Nil(term)

	// length of 5 instead of 4
	raw = []byte{'X', '\x00', '\x00', '\x00', '\x05'}

	term, err = pgproto.ParseTermination(bytes.NewReader(raw))
	s.NotNil(err)
	s.Nil(term)
}

func (s *TerminationTestSuite) Test_Termination_Encode() {
	term := &pgproto.Termination{}
	s.Equal(rawTerminationMessage[:], term.Encode())
}

func BenchmarkTermination_Encode(b *testing.B) {
	term := &pgproto.Termination{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			term.Encode()
		}
	})
}
