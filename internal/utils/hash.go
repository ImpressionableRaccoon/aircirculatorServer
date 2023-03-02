package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func GetMD5(data []byte) (hash string) {
	h := md5.New()
	h.Write(data)
	res := h.Sum(nil)

	hash = hex.EncodeToString(res)
	return hash
}

func GetSHA256(data []byte) (hash string) {
	h := sha256.New()
	h.Write(data)
	res := h.Sum(nil)

	hash = hex.EncodeToString(res)
	return hash
}
