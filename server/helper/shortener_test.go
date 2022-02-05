package helper

import (
	"log"
	"testing"
)

func TestEncode(t *testing.T) {
	var data2Encode uint64 = 7

	encoded_url := Encode(data2Encode)
	expected, decodeErr := Decode(encoded_url)
	if decodeErr != nil {
		log.Fatalln("Error generating QR code. ", decodeErr)
	}

	if data2Encode != expected {
		t.Errorf("got %q, wanted %q", expected, 7)
	}
}
