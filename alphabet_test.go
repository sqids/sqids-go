package sqids

import (
	"reflect"
	"testing"
)

func TestAlphabetSimple(t *testing.T) {
	numbers := []uint64{1, 2, 3}
	id := "4d9fd2"

	alphabet := "0123456789abcdef"
	s, err := NewCustom(Options{
		Alphabet: &alphabet,
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

	decodedNumbers := s.Decode(generatedID)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
	}
}

func TestShortAlphabet(t *testing.T) {
	alphabet := "abcde"
	s, err := NewCustom(Options{
		Alphabet: &alphabet,
	})
	if err != nil {
		t.Fatal(err)
	}

	numbers := []uint64{1, 2, 3}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(generatedID)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", generatedID, numbers, decodedNumbers)
	}
}
