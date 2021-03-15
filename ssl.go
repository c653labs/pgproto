package pgproto

import (
	"io"
)

type SSLRequest struct {
}

func (s *SSLRequest) client() {}

func ParseSSLResponse(r io.Reader) error {
	b := newReadBuffer(r)
	return b.ReadTag('S')
}

func (s *SSLRequest) Encode() []byte {
	// [int32 - length] []byte \0
	//w := newWriteBuffer()
	//w.WriteInt(sslRequestVersion)
	//w.PrependLength()
	return []byte{0x0, 0x0, 0x0, 0x8, 0x4, 0xd2, 0x16, 0x2f}
}

func (s *SSLRequest) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "SSLRequest",
		"Payload": nil,
	}
}

func (s *SSLRequest) String() string { return messageToString(s) }
