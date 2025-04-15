package hdwallet

import (
	"crypto/hmac"
	"crypto/sha512"
)

func eraseBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

func sha512hmac(key, data []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	return h.Sum(nil)
}
