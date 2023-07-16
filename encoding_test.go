package sqids

import (
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	numbers := []uint64{1, 2, 3}

	sqids, err := New()
	if err != nil {
		t.Fatal(err)
	}

	id, err := sqids.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decoded := sqids.Decode(id)

	if !reflect.DeepEqual(numbers, decoded) {
		t.Errorf("Could not encode/decode `%v`", numbers)
	}
}
