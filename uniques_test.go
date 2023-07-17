package sqids

import (
	"reflect"
	"testing"
)

const upper = 1_000_000

func TestUniquesWithPadding(t *testing.T) {
	minLength := len(DefaultAlphabet)
	s, err := NewCustom(Options{
		MinLength: &minLength,
	})
	if err != nil {
		t.Fatal(err)
	}

	set := make(map[string]struct{})

	for i := uint64(0); i < upper; i++ {
		numbers := []uint64{i}
		id, _ := s.Encode(numbers)
		set[id] = struct{}{}

		decodedNumbers := s.Decode(id)
		if !reflect.DeepEqual(numbers, decodedNumbers) {
			t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
		}
	}

	if len(set) != upper {
		t.Errorf("Invalid set count")
	}
}

func TestUniquesLowRanges(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	set := make(map[string]struct{})

	for i := uint64(0); i < upper; i++ {
		numbers := []uint64{i}
		id, _ := s.Encode(numbers)
		set[id] = struct{}{}

		decodedNumbers := s.Decode(id)
		if !reflect.DeepEqual(numbers, decodedNumbers) {
			t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
		}
	}

	if len(set) != upper {
		t.Errorf("Invalid set count")
	}
}

func TestUniquesHighRanges(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	set := make(map[string]struct{})

	for i := uint64(100_000_000); i < 100_000_000+upper; i++ {
		numbers := []uint64{i}
		id, _ := s.Encode(numbers)
		set[id] = struct{}{}

		decodedNumbers := s.Decode(id)
		if !reflect.DeepEqual(numbers, decodedNumbers) {
			t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
		}
	}

	if len(set) != upper {
		t.Errorf("Invalid set count")
	}
}

func TestUniquesMulti(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	set := make(map[string]struct{})

	for i := uint64(0); i < upper; i++ {
		numbers := []uint64{i, i, i, i, i}
		id, _ := s.Encode(numbers)
		set[id] = struct{}{}

		decodedNumbers := s.Decode(id)
		if !reflect.DeepEqual(numbers, decodedNumbers) {
			t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
		}
	}

	if len(set) != upper {
		t.Errorf("Invalid set count")
	}
}
