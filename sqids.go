package sqids

//go:generate go run github.com/campoy/embedmd/v2@v2.0.0 -w README.md

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	defaultAlphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	minAlphabetLength = 5
	minUint64Value    = uint64(0)
	maxUint64Value    = uint64(math.MaxUint64)
)

var defaultBlocklist []string = newDefaultBlocklist()

// Alphabet validation errors
var (
	errAlphabetMultibyte      = errors.New("alphabet must not contain any multibyte characters")
	errAlphabetTooShort       = errors.New("alphabet length must be at least 5")
	errAlphabetNotUniqueChars = errors.New("alphabet must contain unique characters")
	errAlphabetMinLength      = errors.New("alphabet minimum length")
)

// Options for a custom instance of Sqids
type Options struct {
	Alphabet  string
	MinLength int
	Blocklist []string
}

// Sqids lets you generate unique IDs from numbers
type Sqids struct {
	alphabet  string
	minLength int
	blocklist []string
}

// New constructs an instance of Sqids
func New(options ...Options) (*Sqids, error) {
	if len(options) == 0 {
		options = append(options, Options{
			Alphabet:  defaultAlphabet,
			Blocklist: defaultBlocklist,
		})
	}

	// Validate the first given options value, or the default options if none were given.
	o, err := validatedOptions(options[0])
	if err != nil {
		return nil, err
	}

	return &Sqids{
		alphabet:  shuffle(o.Alphabet),
		minLength: o.MinLength,
		blocklist: o.Blocklist,
	}, nil
}

func validatedOptions(o Options) (Options, error) {
	if o.Alphabet == "" {
		o.Alphabet = defaultAlphabet
	}

	// check that the alphabet does not contain multibyte characters
	if len(o.Alphabet) != len([]rune(o.Alphabet)) {
		return Options{}, errAlphabetMultibyte
	}

	// check the length of the alphabet
	if len(o.Alphabet) < minAlphabetLength {
		return Options{}, errAlphabetTooShort
	}

	// check that the alphabet has only unique characters
	if !hasUniqueChars(o.Alphabet) {
		return Options{}, errAlphabetNotUniqueChars
	}

	// test min length (type [might be lang-specific] + min length + max length)
	if o.MinLength < int(minUint64Value) || o.MinLength > len(o.Alphabet) {
		return Options{}, fmt.Errorf("%w has to be between %d and %d", errAlphabetMinLength, minUint64Value, len(o.Alphabet))
	}

	o.Blocklist = filterBlocklist(o.Alphabet, o.Blocklist)

	return o, nil
}

// Encode a slice of uint64 values into an ID string
func (s *Sqids) Encode(numbers []uint64) (string, error) {
	// if no numbers passed, return an empty string
	if len(numbers) == 0 {
		return "", nil
	}

	return s.encodeNumbers(numbers, false)
}

func (s *Sqids) encodeNumbers(numbers []uint64, partitioned bool) (string, error) {
	var (
		err       error
		offset    = calculateOffset(s.alphabet, numbers)
		alphabet  = alphabetOffset(s.alphabet, offset)
		prefix    = alphabet[0]
		partition = alphabet[1]
		ret       = []rune{prefix}
	)

	alphabet = alphabet[2:]

	for i, num := range numbers {
		alphabetWithoutSeparator := alphabet[:len(alphabet)-1]

		ret = append(ret, []rune(toID(num, string(alphabetWithoutSeparator)))...)

		if i < len(numbers)-1 {
			var separator rune

			if partitioned && i == 0 {
				separator = partition
			} else {
				separator = alphabet[len(alphabet)-1]
			}

			ret = append(ret, separator)

			alphabet = []rune(shuffle(string(alphabet)))
		}
	}

	id := string(ret)

	if s.minLength > len(id) {
		if !partitioned {
			numbers = append([]uint64{0}, numbers...)

			id, err = s.encodeNumbers(numbers, true)
			if err != nil {
				return "", err
			}
		}

		if s.minLength > len(id) {
			id = id[:1] + string(alphabet[:s.minLength-len(id)]) + id[1:]
		}
	}

	if s.isBlockedID(id) {
		if partitioned {
			if numbers[0] == maxUint64Value {
				return "", errors.New("ran out of range checking against the blocklist")
			}

			numbers[0]++
		} else {
			numbers = append([]uint64{0}, numbers...)
		}

		id, err = s.encodeNumbers(numbers, true)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

// Decode id string into a slice of uint64 values
func (s *Sqids) Decode(id string) []uint64 {
	ret := []uint64{}

	if id == "" {
		return ret
	}

	rid := []rune(id)

	alphabet := []rune(s.alphabet)

	for _, r := range rid {
		if !contains(alphabet, r) {
			return ret
		}
	}

	prefix := rid[0]
	offset := index(alphabet, prefix)

	alphabet = alphabetOffset(s.alphabet, offset)

	partition := alphabet[1]

	rid = rid[1:]
	alphabet = alphabet[2:]

	if pi := index(rid, partition); pi > 0 && pi < len(rid)-1 {
		rid = rid[pi+1:]
		alphabet = shuffleRunes(alphabet)
	}

	for len(rid) > 0 {
		separator := alphabet[len(alphabet)-1]

		chunks := splitChunks(rid, separator)

		if len(chunks) > 0 {
			alphabetWithoutSeparator := alphabet[:len(alphabet)-1]
			charSet := make(map[rune]bool)

			for _, c := range alphabetWithoutSeparator {
				charSet[c] = true
			}

			for _, c := range chunks[0] {
				if _, exists := charSet[c]; !exists {
					return []uint64{}
				}
			}

			ret = append(ret, toNumber(chunks[0], alphabetWithoutSeparator))

			if len(chunks) > 1 {
				alphabet = shuffleRunes(alphabet)
			}
		}

		if len(chunks) > 0 {
			rid = joinRuneSlices(chunks[1:], separator)
		} else {
			return []uint64{}
		}
	}

	return ret
}

func alphabetOffset(alphabet string, offset int) []rune {
	runes := []rune(alphabet)

	return append(runes[offset:], runes[:offset]...)
}

func joinRuneSlices(rs [][]rune, separator rune) []rune {
	var runes []rune

	if len(rs) > 0 {
		for _, s := range rs[:len(rs)-1] {
			runes = append(runes, s...)
			runes = append(runes, separator)
		}

		runes = append(runes, rs[len(rs)-1]...)
	}

	return runes
}

func splitChunks(runes []rune, separator rune) [][]rune {
	var n int

	var out [][]rune

	for _, r := range runes {
		if r == separator {
			n++
		}

		if len(out) == n {
			out = append(out, []rune{})
		}

		if r != separator {
			out[n] = append(out[n], r)
		}
	}

	return out
}

func (s *Sqids) isBlockedID(id string) bool {
	id = strings.ToLower(id)

	for _, word := range s.blocklist {
		if len(word) <= len(id) {
			if len(id) <= 3 || len(word) <= 3 {
				if id == word {
					return true
				}
			} else if hasDigit(word) {
				if strings.HasPrefix(id, word) || strings.HasSuffix(id, word) {
					return true
				}
			} else if strings.Contains(id, word) {
				return true
			}
		}
	}

	return false
}

// MinValue returns the minimum uint64 value, which is 0
func MinValue() uint64 {
	return minUint64Value
}

// MaxValue returns the maximum uint64 value, which is 18446744073709551615
func MaxValue() uint64 {
	return maxUint64Value
}

func calculateOffset(alphabet string, numbers []uint64) int {
	var (
		offset = len(numbers)
		runes  = []rune(alphabet)
		count  = uint64(len(runes))
	)

	if offset == 0 || len(runes) == 0 {
		return -1
	}

	for i, v := range numbers {
		offset += int(runes[v%count]) + i
	}

	return offset % len(runes)
}

func shuffle(alphabet string) string {
	return string(shuffleRunes([]rune(alphabet)))
}

func shuffleRunes(runes []rune) []rune {
	for i, j := 0, len(runes)-1; j > 0; i, j = i+1, j-1 {
		r := (i*j + int(runes[i]) + int(runes[j])) % len(runes)
		runes[i], runes[r] = runes[r], runes[i]
	}

	return runes
}

func toID(num uint64, alphabet string) string {
	var (
		id     = []rune{}
		runes  = []rune(alphabet)
		count  = uint64(len(runes))
		result = num
	)

	for {
		index := result % count

		id = append([]rune{runes[index]}, id...)

		result = result / count

		if result == 0 {
			break
		}
	}

	return string(id)
}

func toNumber(rid []rune, runes []rune) uint64 {
	count := uint64(len(runes))

	var result uint64

	for _, r := range rid {
		result = (result * count) + uint64(index(runes, r))
	}

	return result
}

func index(s []rune, r rune) int {
	for i := range s {
		if r == s[i] {
			return i
		}
	}

	return -1
}

func hasUniqueChars(str string) bool {
	charSet := make(map[rune]bool)
	for _, c := range str {
		if _, ok := charSet[c]; ok {
			return false
		}
		charSet[c] = true
	}
	return true
}

func contains(s []rune, r rune) bool {
	for _, v := range s {
		if v == r {
			return true
		}
	}

	return false
}

func hasDigit(word string) bool {
	for _, r := range word {
		if r >= '0' && r <= '9' {
			return true
		}
	}

	return false
}
