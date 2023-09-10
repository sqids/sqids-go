package sqids

import "testing"

func TestCalculateOffset(t *testing.T) {
	for _, tt := range []struct {
		alphabet string
		numbers  []uint64
		want     int
	}{
		{"", []uint64{}, -1},
		{"", []uint64{0}, -1},
		{"abcde", []uint64{0}, 3},
		{"fghij", []uint64{0}, 3},
		{"abcde", []uint64{1}, 4},
		{"abcde", []uint64{2}, 0},
		{defaultAlphabet, []uint64{24}, 60},
		{defaultAlphabet, []uint64{25}, 61},
		{defaultAlphabet, []uint64{26}, 4},
		{defaultAlphabet, []uint64{27}, 5},
		{defaultAlphabet, []uint64{1, 2, 3}, 55},
		{defaultAlphabet, []uint64{4, 5, 6}, 2},
	} {
		if _, err := New(Options{Alphabet: tt.alphabet}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := calculateOffset(tt.alphabet, tt.numbers, 0); got != tt.want {
			t.Fatalf("calculateOffset(%q, %#v) = %d, want %d", tt.alphabet, tt.numbers, got, tt.want)
		}
	}
}

func TestDecode(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := len(s.Decode("Re")), 0; got != want {
		t.Fatalf(`len(s.Decode("Re")) = %d, want %d`, got, want)
	}

	if got, want := len(s.Decode("U9")), 1; got != want {
		t.Fatalf(`len(s.Decode("U9")) = %d, want %d`, got, want)
	}
}
