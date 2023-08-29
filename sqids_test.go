package sqids

import "testing"

func TestMinValue(t *testing.T) {
	if got, want := MinValue(), uint64(0); got != want {
		t.Fatalf("MinValue() = %d, want %d", got, want)
	}
}

func TestMaxValue(t *testing.T) {
	if got, want := MaxValue(), uint64(18446744073709551615); got != want {
		t.Fatalf("MaxValue() = %d, want %d", got, want)
	}
}
