package utils

import (
	"encoding/base32"
	"fmt"
)

type Base32Enc struct {
	enc *base32.Encoding
}

func NewBase32Enc(baseChars string) (*Base32Enc, error) {
	if len(baseChars) == 0 {
		return &Base32Enc{base32.StdEncoding.WithPadding(base32.NoPadding)}, nil
	}
	if len(baseChars) != 32 {
		return nil, fmt.Errorf("size of baseChars must be 32")
	}
	return &Base32Enc{base32.NewEncoding(baseChars).WithPadding(base32.NoPadding)}, nil
}

func (b32 *Base32Enc) EncodeToString(b []byte) string {
	return b32.enc.EncodeToString(b)
}

func (b32 *Base32Enc) EncodeStr(s string) string {
	return b32.enc.EncodeToString([]byte(s))
}

func (b32 *Base32Enc) Decode(s string) ([]byte, error) {
	return b32.enc.DecodeString(s)
}
