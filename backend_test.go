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

func (s *BackendKeyDataTestSuite) Test_ParseBackendKeyData() {
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

func BenchmarkBackendKeyDataParse(b *testing.B) {
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

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := pgproto.ParseBackendKeyData(bytes.NewReader(raw))
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func (s *BackendKeyDataTestSuite) Test_ParseBackendKeyData_Empty() {
	backend, err := pgproto.ParseBackendKeyData(bytes.NewReader([]byte{}))
	s.NotNil(err)
	s.Nil(backend)
}

func BenchmarkBackendKeyDataParse_Empty(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pgproto.ParseBackendKeyData(bytes.NewReader([]byte{}))
		}
	})
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
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			m.Encode()
		}
	})
}

func (s *BackendKeyDataTestSuite) Test_BackendKeyData_ParseServerRequest() {
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

	m, err := pgproto.ParseServerMessage(bytes.NewReader(raw))
	backend, ok := m.(*pgproto.BackendKeyData)
	s.Nil(err)
	s.True(ok)
	s.NotNil(backend)
	s.Equal(backend.PID, 1234)
	s.Equal(backend.Key, 1234)
	s.Equal(raw, backend.Encode())
	s.Equal(raw, m.Encode())
}

func BenchmarkBackendKeyData_ParseServerMessage(b *testing.B) {
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

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, err := pgproto.ParseServerMessage(bytes.NewReader(raw))
			if err != nil {
				b.Error(err)
			}
		}
	})
}
