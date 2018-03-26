package pgproto

import (
	"fmt"
	"io"
)

type AuthenticationMethod int

const (
	AUTHENTICATION_OK        AuthenticationMethod = 0
	AUTHENTICATION_PLAINTEXT AuthenticationMethod = 1
	AUTHENTICATION_MD5       AuthenticationMethod = 5
)

func (a AuthenticationMethod) String() string {
	switch a {
	case AUTHENTICATION_OK:
		return "OK"
	case AUTHENTICATION_PLAINTEXT:
		return "Plaintext"
	case AUTHENTICATION_MD5:
		return "MD5"
	}

	return "Unknown"
}

type AuthenticationRequest struct {
	Method AuthenticationMethod
	Salt   []byte
}

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
	if m != AUTHENTICATION_OK && m != AUTHENTICATION_PLAINTEXT && m != AUTHENTICATION_MD5 {
		return nil, fmt.Errorf("received unknown authentication request method number %d", m)
	}

	a := &AuthenticationRequest{
		Method: m,
		Salt:   nil,
	}

	if a.Method == AUTHENTICATION_MD5 {
		a.Salt, err = buf.ReadString(false)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *AuthenticationRequest) Encode() []byte {
	// 'R' [int32 - length] [int32 - method] [other - optional]
	w := newWriteBuffer()
	w.WriteInt(int(a.Method))
	if a.Method == AUTHENTICATION_MD5 {
		w.WriteString(a.Salt, false)
	}
	w.Wrap('R')
	return w.Bytes()
}

func (a *AuthenticationRequest) WriteTo(w io.Writer) (int64, error) { return writeTo(a, w) }

func (a *AuthenticationRequest) String() string {
	return fmt.Sprintf("AuthenticationRequest<Method=%v, Salt=%v>", a.Method, a.Salt)
}
