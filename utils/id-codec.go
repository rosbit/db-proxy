package utils

import (
	"encoding/binary"
	"fmt"
)

const (
	ID_CONN uint8 = iota + 1
	ID_TX
	ID_STMT
)

var idB32 *Base32Enc

func InitIdCodec(b32BaseChars string) (err error) {
	idB32, err = NewBase32Enc(b32BaseChars)
	return
}

func EncodeId(idType uint8, id uint32) string {
	b := make([]byte, 5)
	binary.BigEndian.PutUint32(b, id)
	b[4] = idType
	return idB32.EncodeToString(b)
}

func DecodeId(xId string) (id uint32, idType uint8, err error) {
	b, e := idB32.Decode(xId)
	if e != nil {
		err = e
		return
	}
	idType = b[4]
	switch idType {
	default:
		err = fmt.Errorf("unknown idType")
		return
	case ID_CONN, ID_TX, ID_STMT:
	}
	id = binary.BigEndian.Uint32(b)
	return
}

func EncodeStr(s string) string {
	return idB32.EncodeStr(s)
}

func DecodeStr(x string) (s string, err error) {
	b, e := idB32.Decode(x)
	if e != nil {
		err = e
		return
	}
	s = string(b)
	return
}
