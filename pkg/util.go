package pkg

import (
	"crypto/sha256"
	"fmt"
	"hash/fnv"
)

func GetSHA256Hash(plainText string) string {
	h := sha256.New()
	h.Write([]byte(plainText))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetFNVHash(plainText string) string {
	h := fnv.New64a()
	h.Write([]byte(plainText))
	return fmt.Sprintf("%x", h.Sum(nil))
}
