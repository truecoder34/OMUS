package helper

import (
	"errors"
	"math"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = uint64(len(alphabet))
)

func Encode(number uint64) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)
	for ; number > 0; number = number / length {
		encodedBuilder.WriteByte(alphabet[(number % length)])
	}
	return encodedBuilder.String()
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, symbol := range encoded {
		positionInAlpabet := strings.IndexRune(alphabet, symbol)

		if positionInAlpabet == -1 {
			return uint64(positionInAlpabet), errors.New("invalid character: " + string(symbol))
		}

		number += uint64(positionInAlpabet) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}
