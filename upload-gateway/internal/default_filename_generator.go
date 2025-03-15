package internal

import (
	"crypto/rand"
	"math/big"
)

type DefaultFileNameGenerator struct {
}

func NewDefaultFileNameNameGenerator() *DefaultFileNameGenerator {
	return &DefaultFileNameGenerator{}
}

func (f *DefaultFileNameGenerator) Generate() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 100
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
