package pgproto_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgproto"
	"github.com/stretchr/testify/suite"
)

type BackendKeyDataTestSuite struct {
	suite.Suite
}

func TestBackendKeyDataTestSuite(t *testing.T) {
	suite.Run(t, new(BackendKeyDataTestSuite))
}

func (s *BackendKeyDataTestSuite) Test_ParseBackendKeyData_MD5() {
	raw := []byte{
		// Tag
		'K',
		// Length
		'\x00', '\x00', '\x00', '\x0c',
		// PID
		'\x00', '\x00', '\x04', '\xd2',
		// Key
		'\x00', '\x00', '\x04', '\xd2',
	}

	backend, err := pgproto.ParseBackendKeyData(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(backend)
	s.Equal(backend.PID, 1234)
	s.Equal(backend.Key, 1234)
	s.Equal(raw, backend.Encode())
}

func BenchmarkBackendKeyData_MD5(b *testing.B) {
	raw := []byte{
		// Tag
		'K',
		// Length
		'\x00', '\x00', '\x00', '\x0c',
		// PID
		'\x00', '\x00', '\x04', '\xd2',
		// Key
		'\x00', '\x00', '\x04', '\xd2',
	}

	for i := 0; i < b.N; i++ {
		_, err := pgproto.ParseBackendKeyData(bytes.NewReader(raw))
		if err != nil {
			b.Error(err)
		}
	}
}

func (s *BackendKeyDataTestSuite) Test_ParseBackendKeyData_Empty() {
	backend, err := pgproto.ParseBackendKeyData(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(backend)
}

func BenchmarkBackendKeyData_Empty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pgproto.ParseBackendKeyData(bytes.NewReader([]byte{}))
	}
}

func (s *BackendKeyDataTestSuite) Test_BackendKeyDataEncode() {
	expected := []byte{
		// Tag
		'K',
		// Length
		'\x00', '\x00', '\x00', '\x0c',
		// PID
		'\x00', '\x00', '\x04', '\xd2',
		// Key
		'\x00', '\x00', '\x04', '\xd2',
	}

	b := &pgproto.BackendKeyData{
		PID: 1234,
		Key: 1234,
	}
	s.Equal(expected, b.Encode())
}

func BenchmarkBackendKeyDataEncode(b *testing.B) {
	m := &pgproto.BackendKeyData{
		PID: 1234,
		Key: 1234,
	}
	for i := 0; i < b.N; i++ {
		m.Encode()
	}
}
