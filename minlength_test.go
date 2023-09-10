package sqids

import (
	"reflect"
	"testing"
)

func TestMinLengthSimple(t *testing.T) {
	s, err := New(Options{
		MinLength: uint8(len(defaultAlphabet)),
	})
	if err != nil {
		t.Fatal(err)
	}

	numbers := []uint64{1, 2, 3}

	id := "86Rf07xd4zBmiJXQG6otHEbew02c3PWsUOLZxADhCpKj7aVFv9I8RquYrNlSTM"

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

func TestMinLengthIncremental(t *testing.T) {
	numbers := []uint64{1, 2, 3}
	m := map[uint8]string{
		6:  "86Rf07",
		7:  "86Rf07x",
		8:  "86Rf07xd",
		9:  "86Rf07xd4",
		10: "86Rf07xd4z",
		11: "86Rf07xd4zB",
		12: "86Rf07xd4zBm",
		13: "86Rf07xd4zBmi",
	}
	m[uint8(len(defaultAlphabet))+0] = "86Rf07xd4zBmiJXQG6otHEbew02c3PWsUOLZxADhCpKj7aVFv9I8RquYrNlSTM"
	m[uint8(len(defaultAlphabet))+1] = "86Rf07xd4zBmiJXQG6otHEbew02c3PWsUOLZxADhCpKj7aVFv9I8RquYrNlSTMy"
	m[uint8(len(defaultAlphabet))+2] = "86Rf07xd4zBmiJXQG6otHEbew02c3PWsUOLZxADhCpKj7aVFv9I8RquYrNlSTMyf"
	m[uint8(len(defaultAlphabet))+3] = "86Rf07xd4zBmiJXQG6otHEbew02c3PWsUOLZxADhCpKj7aVFv9I8RquYrNlSTMyf1"

	for minLength, id := range m {
		s, err := New(Options{
			MinLength: minLength,
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

		if len(generatedID) != int(minLength) {
			t.Errorf("Encoding `%v` should produce `%v` length, but produced `%v` length instead", numbers, minLength, len(generatedID))
		}

		decodedNumbers := s.Decode(id)
		if !reflect.DeepEqual(numbers, decodedNumbers) {
			t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
		}
	}
}

func TestMinLengthIncrementalNumbers(t *testing.T) {
	s, err := New(Options{
		MinLength: uint8(len(defaultAlphabet)),
	})
	if err != nil {
		t.Fatal(err)
	}

	ids := map[string][]uint64{
		"SvIzsqYMyQwI3GWgJAe17URxX8V924Co0DaTZLtFjHriEn5bPhcSkfmvOslpBu": {0, 0},
		"n3qafPOLKdfHpuNw3M61r95svbeJGk7aAEgYn4WlSjXURmF8IDqZBy0CT2VxQc": {0, 1},
		"tryFJbWcFMiYPg8sASm51uIV93GXTnvRzyfLleh06CpodJD42B7OraKtkQNxUZ": {0, 2},
		"eg6ql0A3XmvPoCzMlB6DraNGcWSIy5VR8iYup2Qk4tjZFKe1hbwfgHdUTsnLqE": {0, 3},
		"rSCFlp0rB2inEljaRdxKt7FkIbODSf8wYgTsZM1HL9JzN35cyoqueUvVWCm4hX": {0, 4},
		"sR8xjC8WQkOwo74PnglH1YFdTI0eaf56RGVSitzbjuZ3shNUXBrqLxEJyAmKv2": {0, 5},
		"uY2MYFqCLpgx5XQcjdtZK286AwWV7IBGEfuS9yTmbJvkzoUPeYRHr4iDs3naN0": {0, 6},
		"74dID7X28VLQhBlnGmjZrec5wTA1fqpWtK4YkaoEIM9SRNiC3gUJH0OFvsPDdy": {0, 7},
		"30WXpesPhgKiEI5RHTY7xbB1GnytJvXOl2p0AcUjdF6waZDo9Qk8VLzMuWrqCS": {0, 8},
		"moxr3HqLAK0GsTND6jowfZz3SUx7cQ8aC54Pl1RbIvFXmEJuBMYVeW9yrdOtin": {0, 9},
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

func TestMinLengths(t *testing.T) {
	for _, minLength := range []uint8{0, 1, 5, 10, uint8(len(defaultAlphabet))} {
		for _, numbers := range [][]uint64{
			{minUint64Value},
			{0, 0, 0, 0, 0},
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{100, 200, 300},
			{1000, 2000, 3000},
			{1000000},
			{maxUint64Value},
		} {
			s, err := New(Options{
				MinLength: minLength,
			})
			if err != nil {
				t.Fatal(err)
			}

			generatedID, err := s.Encode(numbers)
			if err != nil {
				t.Fatal(err)
			}

			if uint8(len(generatedID)) < minLength {
				t.Errorf("Encoding `%v` with min length `%v` produced `%v`", numbers, minLength, generatedID)
			}

			decodedNumbers := s.Decode(generatedID)
			if !reflect.DeepEqual(numbers, decodedNumbers) {
				t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", generatedID, numbers, decodedNumbers)
			}
		}
	}
}
