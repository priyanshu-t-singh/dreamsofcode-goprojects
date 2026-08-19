package url

import (
	"crypto/rand"
	"math/big"
)

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(length int) (string, error) {
	b := make([]byte, length)
	maxIndex := big.NewInt(int64(len(base62Chars)))

	for i := range length {
		n, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "", err
		}
		b[i] = base62Chars[n.Int64()]
	}

	return string(b), nil
}
