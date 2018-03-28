package pgproto

import (
	"fmt"
	"io"
)

// AuthenticationMethod represents the authentication method requested by the server
type AuthenticationMethod int

// Available authentication methods
const (
	AuthenticationMethodOK        AuthenticationMethod = 0
	AuthenticationMethodPlaintext AuthenticationMethod = 3
	AuthenticationMethodMD5       AuthenticationMethod = 5
)

func (a AuthenticationMethod) String() string {
	switch a {
	case AuthenticationMethodOK:
		return "OK"
	case AuthenticationMethodPlaintext:
		return "Plaintext"
	case AuthenticationMethodMD5:
		return "MD5"
	}

	return "Unknown"
}

// AuthenticationRequest is a server response either asking the client to authenticate or
// used to indicate that authentication was successful
type AuthenticationRequest struct {
	Method AuthenticationMethod
	Salt   []byte
}

func (a *AuthenticationRequest) server() {}

// ParseAuthenticationRequest will attempt to read an AuthenticationRequest message from the reader
func ParseAuthenticationRequest(r io.Reader) (*AuthenticationRequest, error) {
	b := newReadBuffer(r)

	// 'R' [int32 - length] [int32 - method] [other - optional]
	err := b.ReadTag('R')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	// Method
	i, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	m := AuthenticationMethod(i)
	if m != AuthenticationMethodOK && m != AuthenticationMethodPlaintext && m != AuthenticationMethodMD5 {
		return nil, fmt.Errorf("received unknown authentication request method number %d", m)
	}

	a := &AuthenticationRequest{
		Method: m,
	}

	if a.Method == AuthenticationMethodMD5 {
		a.Salt, err = buf.ReadString(false)
		if err != nil {
			return nil, err
		}
		if len(a.Salt) != 4 {
			return nil, fmt.Errorf("expected salt of length 4")
		}
	}

	return a, nil
}

// Encode will return the byte representation of this AuthenticationRequest message
func (a *AuthenticationRequest) Encode() []byte {
	// 'R' [int32 - length] [int32 - method] [other - optional]
	w := newWriteBuffer()
	w.WriteInt(int(a.Method))
	if a.Method == AuthenticationMethodMD5 {
		w.WriteString(a.Salt, false)
	}
	w.Wrap('R')
	return w.Bytes()
}

// WriteTo will write the byte representation of this AuthenticationRequest to the provided writer
func (a *AuthenticationRequest) WriteTo(w io.Writer) (int64, error) { return writeTo(a, w) }

func (a *AuthenticationRequest) String() string {
	return fmt.Sprintf("AuthenticationRequest<Method=%v, Salt=%v>", a.Method, a.Salt)
}
