package sqids

import (
	"fmt"
	"testing"
)

func FuzzDecode(f *testing.F) {
	s, err := New()
	if err != nil {
		f.Fatalf("unexpected error: %v", err)
	}

	f.Fuzz(func(t *testing.T, id string) {
		s.Decode(id)
	})
}

func FuzzEncode(f *testing.F) {
	s, err := New()
	if err != nil {
		f.Fatalf("unexpected error: %v", err)
	}

	f.Fuzz(func(t *testing.T, u uint64) {
		s.Encode([]uint64{u})
	})
}

func FuzzNewEncodeDecode(f *testing.F) {
	f.Add(defaultAlphabet, uint64(1))

	f.Fuzz(func(t *testing.T, alphabet string, u uint64) {
		s, err := New(Options{
			Alphabet: alphabet,
		})
		if err == nil {
			id, err := s.Encode([]uint64{u})
			if err == nil {
				d := s.Decode(id)

				if d[0] != u {
					panic(fmt.Sprintf("%d != %d", d[0], u))
				}
			}
		}
	})
}
