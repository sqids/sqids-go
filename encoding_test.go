package sqids

import (
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
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
