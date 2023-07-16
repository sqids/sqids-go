package sqids

import (
	"reflect"
	"testing"
)

func TestBlocklist(t *testing.T) {
	numbers := []uint64{1, 2, 3}
	id := "TM0x1Mxz"

	s, err := NewSqids(Options{
		Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		MinLength: 0,
		Blocklist: &[]string{
			"8QRLaD",   // normal result of 1st encoding, let's block that word on purpose
			"7T1cd0dL", // result of 2nd encoding
			"UeIe",     // result of 3rd encoding is `RA8UeIe7`, let's block a substring
			"imhw",     // result of 4th encoding is `WM3Limhw`, let's block the postfix
			"LfUQ",     // result of 4th encoding is `LfUQh4HN`, let's block the prefix
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	if id != generatedID {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", numbers, id, generatedID)
	}

	decodedNumbers := s.Decode(id)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
	}
}
