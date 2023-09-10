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
	numbers := []uint64{4572721}
	blockedID := "aho1e"
	unblockedID := "JExTR"

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
	numbers := []uint64{4572721}
	id := "aho1e"

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
	numbers := []uint64{4572721}
	id := "aho1e"

	s, err := New(Options{
		Blocklist: []string{
			"ArUO", // originally encoded [100000]
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
	decodedNumbers = s.Decode("ArUO")
	if !reflect.DeepEqual([]uint64{100_000}, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, []uint64{100_000}, decodedNumbers)
	}

	generatedID, err = s.Encode([]uint64{100_000})
	if err != nil {
		t.Fatal(err)
	}

	if generatedID != "QyG4" {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", []uint64{100_000}, "QyG4", generatedID)
	}

	decodedNumbers = s.Decode("QyG4")
	if !reflect.DeepEqual([]uint64{100_000}, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", id, []uint64{100_000}, decodedNumbers)
	}
}

func TestNewBlocklist(t *testing.T) {
	numbers := []uint64{1000000, 2000000}
	id := "1aYeB7bRUt"

	s, err := New(Options{
		Blocklist: []string{
			"JSwXFaosAN", // normal result of 1st encoding, let's block that word on purpose
			"OCjV9JK64o", // result of 2nd encoding
			"rBHf",       // result of 3rd encoding is `4rBHfOiqd3`, let's block a substring
			"79SM",       // result of 4th encoding is `dyhgw479SM`, let's block the postfix
			"7tE6",       // result of 4th encoding is `7tE6jdAHLe`, let's block the prefix
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
	blocklist := []string{"86Rf07", "se8ojk", "ARsz1p", "Q8AI49", "5sQRZO"}

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
		Blocklist: []string{"pnd"},
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
	id := "IBSHOZ"

	s, err := New(Options{
		Alphabet:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Blocklist: []string{"sxnzkl"}, // lowercase blocklist in only-uppercase alphabet
	})
	if err != nil {
		t.Fatal(err)
	}

	generatedID, err := s.Encode(numbers)
	if err != nil {
		t.Fatal(err)
	}

	decodedNumbers := s.Decode(id)

	// without blocklist, would've been "SXNZKL"
	if id != generatedID {
		t.Errorf("Encoding `%v` should produce `%v`, but instead produced `%v`", numbers, id, generatedID)
	}

	if !reflect.DeepEqual(numbers, decodedNumbers) {
		t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", generatedID, numbers, decodedNumbers)
	}
}

func TestMaxEncodingAttempts(t *testing.T) {
	alphabet := "abc"
	minLength := uint8(3)
	blocklist := []string{"cab", "abc", "bca"}

	s, err := New(Options{
		Alphabet:  alphabet,
		MinLength: minLength,
		Blocklist: blocklist,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(alphabet) != int(minLength) || len(blocklist) != int(minLength) {
		t.Errorf("`TestMaxEncodingAttempts` is not setup properly")
	}

	if _, err := s.Encode([]uint64{0}); err == nil {
		t.Errorf("Should throw error about max regeneration attempts")
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
