package sqids

import (
	"reflect"
	"testing"
)

func TestBlocklist(t *testing.T) {
	var (
		defaultBlocklist = Blocklist()
		customBlocklist  = Blocklist("custom1", "custom2")
	)

	if got, want := len(defaultBlocklist), len(newDefaultBlocklist()); got != want {
		t.Fatalf("len(defaultBlocklist) = %d, want %d", got, want)
	}

	if got, want := len(customBlocklist), len(defaultBlocklist)+2; got != want {
		t.Fatalf("len(customBlocklist) = %d, want %d", got, want)
	}
}

func TestBlocklistDefault(t *testing.T) {
	numbers := []uint64{200044}
	blockedID := "sexy"
	unblockedID := "d171vI"

	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(blockedID)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", blockedID, numbers, decodedNumbers)
	}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	if unblockedID != generatedID {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", numbers, unblockedID, generatedID)
	}
}

func TestBlocklistEmpty(t *testing.T) {
	numbers := []uint64{200044}
	id := "sexy"

	s, err := New(Options{
		Blocklist: []string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(id)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
	}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}
	if id != generatedID {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", numbers, id, generatedID)
	}
}

func TestBlocklistNonEmpty(t *testing.T) {
	numbers := []uint64{200044}
	id := "sexy"

	s, err := New(Options{
		Blocklist: []string{
			"AvTg", // originally encoded [100000]
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// make sure we don't use the default blocklist
	decodedNumbers := s.Decode(id)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
	}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	if id != generatedID {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", numbers, id, generatedID)
	}

	// make sure we are using the passed blocklist
	decodedNumbers = s.Decode("AvTg")
	if !reflect.DeepEqual([]uint64{100_000}, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, []uint64{100_000}, decodedNumbers)
	}

	generatedID, err = s.Encode([]uint64{100_000})
	if err != nil {
		t.Fatal(err)
	}

	if generatedID != "7T1X8k" {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", []uint64{100_000}, "7T1X8k", generatedID)
	}

	decodedNumbers = s.Decode("7T1X8k")
	if !reflect.DeepEqual([]uint64{100_000}, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, []uint64{100_000}, decodedNumbers)
	}
}

func TestNewBlocklist(t *testing.T) {
	numbers := []uint64{1, 2, 3}
	id := "TM0x1Mxz"

	s, err := New(Options{
		Blocklist: []string{
			"8QRLaD",   // normal result of 1st encoding, let's block that word on purpose
			"7T1cd0dL", // result of 2nd encoding
			"UeIe",     // result of 3rd encoding is `RA8UeIe7`, let's block a substring
			"imhw",     // result of 4th encoding is `WM3Limhw`, let's block the postfix
			"LfUQ",     // result of 4th encoding is `LfUQh4HN`, let's block the prefix
		},
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

func TestDecodingBlocklistedIDs(t *testing.T) {
	numbers := []uint64{1, 2, 3}
	blocklist := []string{"8QRLaD", "7T1cd0dL", "RA8UeIe7", "WM3Limhw", "LfUQh4HN"}

	s, err := New(Options{
		Blocklist: blocklist,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range blocklist {
		decodedNumbers := s.Decode(id)
		if !reflect.DeepEqual(decodedNumbers, numbers) {
			t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, numbers, decodedNumbers)
		}
	}
}

func TestShortBlocklistMatch(t *testing.T) {
	numbers := []uint64{1_000}

	s, err := New(Options{
		Blocklist: []string{"pPQ"},
	})

	if err != nil {
		t.Fatal(err)
	}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(generatedID)
	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", generatedID, numbers, decodedNumbers)
	}
}

func TestUpperCaseAlphabetBlocklistFiltering(t *testing.T) {
	numbers := []uint64{1, 2, 3}
	id := "ULPBZGBM"

	s, err := New(Options{
		Alphabet:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Blocklist: []string{"sqnmpn"}, // lowercase blocklist in only-uppercase alphabet
	})
	if err != nil {
		t.Fatal(err)
	}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(id)

	// without blocklist, would've been "SQNMPN"
	if id != generatedID {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", numbers, id, generatedID)
	}

	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", generatedID, numbers, decodedNumbers)
	}
}

func TestFilterBlocklist(t *testing.T) {
	t.Run("no words less than 3 chars", func(t *testing.T) {
		filtered := filterBlocklist("YESNO", []string{"yes", "no"})

		if got, want := len(filtered), 1; got != want {
			t.Fatalf("len(filtered) = %d, want %d", got, want)
		}

		if got, want := filtered[0], "yes"; got != want {
			t.Fatalf("filtered[0] = %q, want %q", got, want)
		}
	})

	t.Run("remove words containing letters not in alphabet", func(t *testing.T) {
		filtered := filterBlocklist("YESNO", []string{"yes", "nope"})

		if got, want := len(filtered), 1; got != want {
			t.Fatalf("len(filtered) = %d, want %d", got, want)
		}

		if got, want := filtered[0], "yes"; got != want {
			t.Fatalf("filtered[0] = %q, want %q", got, want)
		}
	})
}
