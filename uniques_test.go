package sqids

import (
	"reflect"
	"testing"
)

func TestUniques(t *testing.T) {
	u := upper()

	t.Run("WithPadding", func(t *testing.T) {
		minLength := len(DefaultAlphabet)

		s, err := NewCustom(Options{
			MinLength: &minLength,
		})
		if err != nil {
			t.Fatal(err)
		}

		set := make(map[string]struct{})

		for i := uint64(0); i < u; i++ {
			numbers := []uint64{i}

			id, _ := s.Encode(numbers)
			set[id] = struct{}{}

			decodedNumbers := s.Decode(id)
			if !reflect.DeepEqual(numbers, decodedNumbers) {
				t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
			}
		}

		if len(set) != int(u) {
			t.Errorf("Invalid set count")
		}
	})

	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("LowRanges", func(t *testing.T) {
		set := make(map[string]struct{})

		for i := uint64(0); i < u; i++ {
			numbers := []uint64{i}

			id, _ := s.Encode(numbers)
			set[id] = struct{}{}

			decodedNumbers := s.Decode(id)
			if !reflect.DeepEqual(numbers, decodedNumbers) {
				t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
			}
		}

		if len(set) != int(u) {
			t.Errorf("Invalid set count")
		}
	})

	t.Run("HighRanges", func(t *testing.T) {
		set := make(map[string]struct{})

		for i := uint64(100_000_000); i < 100_000_000+u; i++ {
			numbers := []uint64{i}

			id, _ := s.Encode(numbers)
			set[id] = struct{}{}

			decodedNumbers := s.Decode(id)
			if !reflect.DeepEqual(numbers, decodedNumbers) {
				t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
			}
		}

		if len(set) != int(u) {
			t.Errorf("Invalid set count")
		}
	})

	t.Run("Multi", func(t *testing.T) {
		set := make(map[string]struct{})

		for i := uint64(0); i < u; i++ {
			numbers := []uint64{i, i, i, i, i}

			id, _ := s.Encode(numbers)
			set[id] = struct{}{}

			decodedNumbers := s.Decode(id)
			if !reflect.DeepEqual(numbers, decodedNumbers) {
				t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
			}
		}

		if len(set) != int(u) {
			t.Errorf("Invalid set count")
		}
	})
}

func upper() uint64 {
	if testing.Short() {
		return 1_000
	}

	return 1_000_000
}
