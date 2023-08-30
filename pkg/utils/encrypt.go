package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(raw []byte) string {
	h := md5.New()
	h.Write(raw)
	x := h.Sum(nil)
	y := make([]byte, 32)
	hex.Encode(y, x)
	return string(y)
}
