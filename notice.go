package pgproto

import "io"

type NoticeResponse Error

func ParseNoticeResponse(r io.Reader) (*NoticeResponse, error) {
	e, err := ParseError(r)
	if err != nil {
		return nil, err
	}

	return (*NoticeResponse)(e), nil
}

func (n *NoticeResponse) server() {}

func (n *NoticeResponse) WriteTo(w io.Writer) (int64, error) { return writeTo(n, w) }

func (n *NoticeResponse) Encode() []byte {
	return encodeError((*Error)(n), 'N')
}

func (n *NoticeResponse) String() string {
	return errorString((*Error)(n), "NoticeResponse")
}
