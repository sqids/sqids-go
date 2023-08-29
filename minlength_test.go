package sqids

import (
	"reflect"
	"testing"
)

func TestMinLengthSimple(t *testing.T) {
	s, err := New(Options{
		MinLength: len(defaultAlphabet),
	})
	if err != nil {
		t.Fatal(err)
	}

	numbers := []uint64{1, 2, 3}

	id := "75JILToVsGerOADWmHlY38xvbaNZKQ9wdFS0B6kcMEtnRpgizhjU42qT1cd0dL"

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

func TestMinLengthIncrementalNumbers(t *testing.T) {
	s, err := New(Options{
		MinLength: len(defaultAlphabet),
	})
	if err != nil {
		t.Fatal(err)
	}

	ids := map[string][]uint64{
		"jf26PLNeO5WbJDUV7FmMtlGXps3CoqkHnZ8cYd19yIiTAQuvKSExzhrRghBlwf": {0, 0},
		"vQLUq7zWXC6k9cNOtgJ2ZK8rbxuipBFAS10yTdYeRa3ojHwGnmMV4PDhESI2jL": {0, 1},
		"YhcpVK3COXbifmnZoLuxWgBQwtjsSaDGAdr0ReTHM16yI9vU8JNzlFq5Eu2oPp": {0, 2},
		"OTkn9daFgDZX6LbmfxI83RSKetJu0APihlsrYoz5pvQw7GyWHEUcN2jBqd4kJ9": {0, 3},
		"h2cV5eLNYj1x4ToZpfM90UlgHBOKikQFvnW36AC8zrmuJ7XdRytIGPawqYEbBe": {0, 4},
		"7Mf0HeUNkpsZOTvmcj836P9EWKaACBubInFJtwXR2DSzgYGhQV5i4lLxoT1qdU": {0, 5},
		"APVSD1ZIY4WGBK75xktMfTev8qsCJw6oyH2j3OnLcXRlhziUmpbuNEar05QCsI": {0, 6},
		"P0LUhnlT76rsWSofOeyRGQZv1cC5qu3dtaJYNEXwk8Vpx92bKiHIz4MgmiDOF7": {0, 7},
		"xAhypZMXYIGCL4uW0te6lsFHaPc3SiD1TBgw5O7bvodzjqUn89JQRfk2Nvm4JI": {0, 8},
		"94dRPIZ6irlXWvTbKywFuAhBoECQOVMjDJp53s2xeqaSzHY8nc17tmkLGwfGNl": {0, 9},
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
	for _, minLength := range []int{0, 1, 5, 10, len(defaultAlphabet)} {
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

			if len(generatedID) < minLength {
				t.Errorf("Encoding `%v` with min length `%v` produced `%v`", numbers, minLength, generatedID)
			}

			decodedNumbers := s.Decode(generatedID)
			if !reflect.DeepEqual(numbers, decodedNumbers) {
				t.Errorf("Decoding `%v` should produce `%v`, but instead produced `%v`", generatedID, numbers, decodedNumbers)
			}
		}
	}
}

func TestOutOfRangeInvalidMinLength(t *testing.T) {
	if _, err := New(Options{
		MinLength: -1,
	}); err == nil {
		t.Errorf("Should not allow out of range min length")
	}

	if _, err := New(Options{
		MinLength: len(defaultAlphabet) + 1,
	}); err == nil {
		t.Errorf("Should not allow out of range min length")
	}
}
