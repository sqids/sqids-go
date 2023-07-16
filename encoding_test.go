package sqids

import (
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	numbers := []uint64{}

	sqids, _ := New()
	id, err := sqids.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := sqids.Decode(id)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(numbers, decoded) {
		t.Errorf("Could not encode/decode `%v`", numbers)
	}
}
