package pgproto

import (
	"fmt"
	"io"
)

// AuthenticationMethod represents the authentication method requested by the server
type AuthenticationMethod int

// Available authentication methods
const (
	AuthenticationMethodOK           AuthenticationMethod = 0
	AuthenticationMethodPlaintext    AuthenticationMethod = 3
	AuthenticationMethodMD5          AuthenticationMethod = 5
	AuthenticationMethodSASL         AuthenticationMethod = 10
	AuthenticationMethodSASLContinue AuthenticationMethod = 11
	AuthenticationMethodSASLFinal    AuthenticationMethod = 12
	SASLMechanismScramSHA256                              = "SCRAM-SHA-256"
	SASLMechanismScramSHA256Plus                          = "SCRAM-SHA-256-PLUS"
)

func (a AuthenticationMethod) String() string {
	switch a {
	case AuthenticationMethodOK:
		return "OK"
	case AuthenticationMethodPlaintext:
		return "Plaintext"
	case AuthenticationMethodMD5:
		return "MD5"
	case AuthenticationMethodSASL:
		return "SASL"
	}

	return "Unknown"
}

// AuthenticationRequest is a server response either asking the client to authenticate or
// used to indicate that authentication was successful
type AuthenticationRequest struct {
	Method                   AuthenticationMethod
	Salt                     []byte
	SupportedScramSHA256     bool
	SupportedScramSHA256Plus bool
	Message                  []byte
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
	switch m := AuthenticationMethod(i); m {
	case AuthenticationMethodOK, AuthenticationMethodPlaintext:
		return &AuthenticationRequest{Method: m}, nil
	case AuthenticationMethodMD5:
		a := &AuthenticationRequest{
			Method: m,
		}
		a.Salt, err = buf.ReadString(false)
		if err != nil {
			return nil, err
		}
		if len(a.Salt) != 4 {
			return nil, fmt.Errorf("expected salt of length 4")
		}
		return a, nil
	case AuthenticationMethodSASL:
		a := &AuthenticationRequest{
			Method: m,
		}

		for {
			supportedSaslMechanism, err := buf.ReadString(true)
			if err != nil {
				return nil, err
			}
			switch string(supportedSaslMechanism) {
			case SASLMechanismScramSHA256:
				a.SupportedScramSHA256 = true
			case SASLMechanismScramSHA256Plus:
				a.SupportedScramSHA256Plus = true
			case "":
				return a, nil
			default:
				return nil, fmt.Errorf("server supports unknown SASL mechanism")
			}
		}
	case AuthenticationMethodSASLContinue, AuthenticationMethodSASLFinal:

		message, err := buf.ReadString(true)
		if err != nil {
			return nil, err
		}
		a := &AuthenticationRequest{
			Method:  m,
			Message: message,
		}

		return a, nil

	default:
		return nil, fmt.Errorf("received unknown authentication request method number %d", m)
	}
}

// Encode will return the byte representation of this AuthenticationRequest message
func (a *AuthenticationRequest) Encode() []byte {
	// 'R' [int32 - length] [int32 - method] [other - optional]
	w := newWriteBuffer()
	w.WriteInt(int(a.Method))
	switch a.Method {
	case AuthenticationMethodMD5:
		w.WriteString(a.Salt, false)
	case AuthenticationMethodSASL:
		if a.SupportedScramSHA256 {
			w.WriteString([]byte(SASLMechanismScramSHA256), true)
		}
		if a.SupportedScramSHA256Plus {
			w.WriteString([]byte(SASLMechanismScramSHA256Plus), true)
		}
		w.WriteByte('\x00')
	case AuthenticationMethodSASLContinue, AuthenticationMethodSASLFinal, AuthenticationMethodOK:
		w.WriteBytes(a.Message)
	}

	w.Wrap('R')
	return w.Bytes()
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "AuthenticationRequest",
//     "Payload": map[string]interface{}{
//       "Method": <AuthenticationRequest.Method>,
//       "Salt": <AuthenticationRequest.Salt>,
//       "ScramSHA256": <AuthenticationRequest.ScramSHA256>,
//       "ScramSHA256Plus": <AuthenticationRequest.ScramSHA256Plus>,
//       "Message": <AuthenticationRequest.Message>,
//     },
//   }
func (a *AuthenticationRequest) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "AuthenticationRequest",
		"Payload": map[string]interface{}{
			"Method":          int(a.Method),
			"Salt":            a.Salt,
			"ScramSHA256":     a.SupportedScramSHA256,
			"ScramSHA256Plus": a.SupportedScramSHA256Plus,
			"Message":         a.Message,
		},
	}
}

func (a *AuthenticationRequest) String() string { return messageToString(a) }
