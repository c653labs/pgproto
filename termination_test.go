package pgmsg_test

import (
	"bytes"
	"testing"

	"github.com/c653labs/pgmsg"
)

func Test_ParseTermination(t *testing.T) {
	raw := []byte{'X', '\x00', '\x00', '\x00', '\x04'}

	term, err := pgmsg.ParseTermination(bytes.NewReader(raw))
	if err != nil {
		t.Errorf("expected err to be nil, instead got %#v", err)
	}

	if !bytes.Equal(raw, term.Encode()) {
		t.Errorf("expected Termination.Encode() to be the same as it's input")
	}
}

func Test_ParseTermination_Empty(t *testing.T) {
	term, err := pgmsg.ParseTermination(bytes.NewReader([]byte{}))
	if err == nil {
		t.Errorf("expected err to not be nil, instead got %#v", err)
	}

	if term != nil {
		t.Errorf("expected term to be nil, instead got %#v", term)
	}
}

func Test_ParseTermination_InvalidTag(t *testing.T) {
	// lowercase 'x' instead of uppercase 'X'
	raw := []byte{'x', '\x00', '\x00', '\x00', '\x04'}

	term, err := pgmsg.ParseTermination(bytes.NewReader(raw))
	if err == nil {
		t.Errorf("expected err to not be nil, instead got %#v", err)
	}

	if term != nil {
		t.Errorf("expected term to be nil, instead got %#v", term)
	}
}

func Test_ParseTermination_InvalidLength(t *testing.T) {
	// length of 3 instead of 4
	raw := []byte{'X', '\x00', '\x00', '\x00', '\x03'}

	term, err := pgmsg.ParseTermination(bytes.NewReader(raw))
	if err == nil {
		t.Errorf("expected err to not be nil, instead got %#v", err)
	}

	if term != nil {
		t.Errorf("expected term to be nil, instead got %#v", term)
	}

	// length of 5 instead of 4
	raw = []byte{'X', '\x00', '\x00', '\x00', '\x05'}

	term, err = pgmsg.ParseTermination(bytes.NewReader(raw))
	if err == nil {
		t.Errorf("expected err to not be nil, instead got %#v", err)
	}

	if term != nil {
		t.Errorf("expected term to be nil, instead got %#v", term)
	}
}

func Test_Termination_Encode(t *testing.T) {
	expected := []byte{'X', '\x00', '\x00', '\x00', '\x04'}

	term := &pgmsg.Termination{}
	if !bytes.Equal(expected, term.Encode()) {
		t.Errorf("expected Termination.Encode() to be the same as it's input")
	}
}
