package helper

import (
	"log"
	"testing"
)

type encodeTest struct {
	arg      uint64
	expected string
}

var encodeTests = []encodeTest{
	{7, "h"},
	{555555, "JGuc"},
	{10000000000, "KY8U4k"},
	{0, ""},
}

type decodeTest struct {
	arg      string
	expected uint64
}

var decodeTests = []decodeTest{
	{"QQQ", 164094},
	{"12Z76", 871288285},
	{"godevblogusinggomodules", 6156118248968608582},
}

func TestEncode(t *testing.T) {
	for idx, testData := range encodeTests {
		encoded_number := Encode(testData.arg)
		if encoded_number != testData.expected {
			t.Errorf("[ TEST %d ] : Output %q not equal to expected %q", idx, encoded_number, testData.expected)
		}
	}
}

func TestDecode(t *testing.T) {
	for idx, testData := range decodeTests {
		number, decodeErr := Decode(testData.arg)
		if decodeErr != nil {
			log.Fatalln("[ERROR] Decoder() error : ", decodeErr)
		}

		if number != testData.expected {
			t.Errorf("[ TEST %d ] : Output %q not equal to expected %q", idx, number, testData.expected)
		}
	}
}
