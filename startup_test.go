package pgmsg_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgmsg"
)

func Test_ParseStartupMessage(t *testing.T) {
}

func Test_ParseStartupMessage_NoOptions(t *testing.T) {
	raw := []byte{'\x00', '\x00', '\x00', '\x08', '\x00', '\x03', '\x00', '\x00'}

	startup, err := pgmsg.ParseStartupMessage(bytes.NewReader(raw))
	if err != nil {
		t.Errorf("expected err to be nil, instead got %#v", err)
	}

	if startup == nil {
		t.Errorf("expected startup to not be nil, instead got %#v", startup)
	}

	if startup.Protocol != 196608 {
		t.Errorf("expected Startup.Protocol to be 196608, instead got %#v", startup.Protocol)
	}
}
