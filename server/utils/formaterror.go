package formaterror

import (
	"errors"
	"strings"
)

/*
	Template for ERRORs handler. need to be extended

*/
func FormatError(err string) error {

	if strings.Contains(err, "OriginalURL") {
		return errors.New("This URL is already encoded")
	}

	if strings.Contains(err, "EncodedURL") {
		return errors.New("This URL is already encoded")
	}

	return errors.New("Incorrect Details")
}
