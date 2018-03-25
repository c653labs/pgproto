package pgmsg_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgmsg"
	"github.com/stretchr/testify/suite"
)

type StartupMessageTestSuite struct {
	suite.Suite
}

func TestStartupMessageTestSuite(t *testing.T) {
	suite.Run(t, new(StartupMessageTestSuite))
}

func (s *StartupMessageTestSuite) Test_ParseStartupMessage() {
	raw := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x25',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// "database" \0
		'\x64', '\x61', '\x74', '\x61', '\x62', '\x61', '\x73', '\x65', '\x00',
		// "db_name" \0
		'\x64', '\x62', '\x5f', '\x6e', '\x61', '\x6d', '\x65', '\x00',
		// "user" \0
		'\x75', '\x73', '\x65', '\x72', '\x00',
		// "pgmsg" \0
		'\x70', '\x67', '\x6d', '\x73', '\x67', '\x00',
		// ending
		'\x00',
	}
	startup, err := pgmsg.ParseStartupMessage(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(startup)
	s.Equal(startup.Protocol, 196608)
	s.Equal(startup.Options["user"], []byte("pgmsg"))
	s.Equal(startup.Options["database"], []byte("db_name"))
	s.Equal(raw, startup.Encode())
}

func (s *StartupMessageTestSuite) Test_ParseStartupMessage_NoOptions() {
	raw := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x09',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// ending
		'\x00',
	}

	startup, err := pgmsg.ParseStartupMessage(bytes.NewReader(raw))
	s.Nil(err)
	s.NotNil(startup)
	s.Equal(startup.Protocol, 196608)
	s.Equal(raw, startup.Encode())
}

func (s *StartupMessageTestSuite) Test_StartupMessageEncode() {
	expected := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x25',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// "database" \0
		'\x64', '\x61', '\x74', '\x61', '\x62', '\x61', '\x73', '\x65', '\x00',
		// "db_name" \0
		'\x64', '\x62', '\x5f', '\x6e', '\x61', '\x6d', '\x65', '\x00',
		// "user" \0
		'\x75', '\x73', '\x65', '\x72', '\x00',
		// "pgmsg" \0
		'\x70', '\x67', '\x6d', '\x73', '\x67', '\x00',
		// ending
		'\x00',
	}
	startup := &pgmsg.StartupMessage{
		Protocol: 196608,
		Options:  make(map[string][]byte),
	}
	startup.Options["user"] = []byte("pgmsg")
	startup.Options["database"] = []byte("db_name")
	s.Equal(expected, startup.Encode())
}
func (s *StartupMessageTestSuite) Test_StartupMessageEncode_NoOptions() {
	expected := []byte{
		// Length
		'\x00', '\x00', '\x00', '\x09',
		// Protocol
		'\x00', '\x03', '\x00', '\x00',
		// ending
		'\x00',
	}
	startup := &pgmsg.StartupMessage{
		Protocol: 196608,
	}
	s.Equal(expected, startup.Encode())
}
