package pgproto

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
)

func bytesToInt(buf []byte) int {
	return int(int32(binary.BigEndian.Uint32(buf)))
}

func bytesToInt16(buf []byte) int {
	return int(int16(binary.BigEndian.Uint16(buf)))
}

func hashPassword(user []byte, password []byte, salt []byte) []byte {
	digest := md5.New()
	digest.Write(password)
	digest.Write(user)
	pwdhash := digest.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(pwdhash)))
	hex.Encode(dst, pwdhash)

	digest = md5.New()
	digest.Write(dst)
	digest.Write(salt)

	hash := digest.Sum(nil)
	dst = make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(dst, hash)

	return append([]byte("md5"), dst...)
}
