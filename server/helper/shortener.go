package helper

import (
	"errors"
	"math"
	"strings"
)

/*
	use 62 possible characters to encode URL. Removins -._~ symbols
*/
const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = uint64(len(alphabet))
)

/*
	encode random number to string. Use formula permutations with repetition ;
	[ param 1 ]: random number in range [ 0 , 18 446 744 073 709 551 615 ]
	[ output ] : encoded string
	TODO: fix case with 0 number. It will return EMPTY string
*/
func Encode(number uint64) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)
	for ; number > 0; number = number / length {
		encodedBuilder.WriteByte(alphabet[(number % length)])
	}
	return encodedBuilder.String()
}

/*
	decode string to number
	[ param 1 ] : encoded string
	[ output 1 ] : initial number,
	[ output 2 ] : error
*/
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
