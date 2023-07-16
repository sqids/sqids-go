package sqids

import (
	"reflect"
	"testing"
)

func TestSimple(t *testing.T) {
	numbers := []uint64{1, 2, 3}
	id := "4d9fd2"

	s, err := NewCustom("0123456789abcdef", 0)
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

func TestShortAlphabet(t *testing.T) {
	s, err := NewCustom("abcde", 0)
	if err != nil {
		t.Fatal(err)
	}

	numbers := []uint64{1, 2, 3}

	id, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(id)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
	}
}
