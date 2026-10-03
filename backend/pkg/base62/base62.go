package base62

import (
	"errors"
	"strings"
)

const (
	aplhabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base     = 62
)

var ErrInvalidInput = errors.New("invalid input for base62 encoding")

func Encode(num int64) (string, error) {
	if num < 0 {
		return "", ErrInvalidInput
	}

	if num == 0 {
		return string(aplhabet[0]), nil
	}

	var result strings.Builder

	for num > 0 {
		remainder := num % base
		result.WriteByte(aplhabet[remainder])
		num /= base
	}

	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}

func Decode(str string) (int64, error) {
	if str == "" {
		return 0, ErrInvalidInput
	}

	var result int64

	for _, char := range str {
		index := strings.IndexRune(aplhabet, char)
		if index == -1 {
			return 0, ErrInvalidInput
		}
		result = result*base + int64(index)
	}

	return result, nil
}
