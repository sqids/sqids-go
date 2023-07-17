package sqids

import (
	"reflect"
	"testing"
)

func TestEncodingSimple(t *testing.T) {
	numbers := []uint64{1, 2, 3}

	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	id, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(id)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Could not encode/decode `%v`", numbers)
	}
}

func TestEncodingDifferentInputs(t *testing.T) {
	numbers := []uint64{0, 0, 0, 1, 2, 3, 100, 1_000, 100_000, 1_000_000, MaxValue()}

	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	id, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(id)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Could not encode/decode `%v`", numbers)
	}
}

func TestEncodingIncrementalNumbers(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	ids := map[string][]uint64{
		"bV": {0},
		"U9": {1},
		"g8": {2},
		"Ez": {3},
		"V8": {4},
		"ul": {5},
		"O3": {6},
		"AF": {7},
		"ph": {8},
		"n8": {9},
	}

	for id, numbers := range ids {
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
}

func TestEncodingIncrementalNumbersSameIndex0(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	ids := map[string][]uint64{
		"SrIu": {0, 0},
		"nZqE": {0, 1},
		"tJyf": {0, 2},
		"e86S": {0, 3},
		"rtC7": {0, 4},
		"sQ8R": {0, 5},
		"uz2n": {0, 6},
		"7Td9": {0, 7},
		"3nWE": {0, 8},
		"mIxM": {0, 9},
	}

	for id, numbers := range ids {
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
}

func TestEncodingIncrementalNumbersSameIndex1(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	ids := map[string][]uint64{
		"SrIu": {0, 0},
		"nbqh": {1, 0},
		"t4yj": {2, 0},
		"eQ6L": {3, 0},
		"r4Cc": {4, 0},
		"sL82": {5, 0},
		"uo2f": {6, 0},
		"7Zdq": {7, 0},
		"36Wf": {8, 0},
		"m4xT": {9, 0},
	}

	for id, numbers := range ids {
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
}

func TestEncodingMultiInput(t *testing.T) {
	numbers := []uint64{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
		26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73,
		74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97,
		98, 99}

	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	id, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(id)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Could not encode/decode `%v`", numbers)
	}
}

func TestEncodingEmptySlice(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	id, err := s.Encode([]uint64{})
	if err != nil {
		t.Fatal(err)
	}

	if id != "" {
		t.Errorf("Could not encode empty slice")
	}
}

func TestEncodingEmptyString(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s.Decode(""), []uint64{}) {
		t.Errorf("Could not decode empty string")
	}
}

func TestEncodingInvalidCharacter(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s.Decode("*"), []uint64{}) {
		t.Errorf("Could not decode with invalid character")
	}
}

// TestEncodingOutOfRange - no need since `[]uint64` handles ranges
// func TestEncodingOutOfRange(t *testing.T) {}
