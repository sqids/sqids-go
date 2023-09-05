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

// Encode -
func (s *Sqids) Encode(numbers []uint64) (string, error) {
	// if no numbers passed, return an empty string
	if len(numbers) == 0 {
		return "", nil
	}

	return s.encodeNumbers(numbers, false)
}

func (s *Sqids) encodeNumbers(numbers []uint64, partitioned bool) (string, error) {
	var err error

	offset := len(numbers)
	for i, v := range numbers {
		offset += int(s.alphabet[v%uint64(len(s.alphabet))]) + i
	}
	offset = offset % len(s.alphabet)

	alphabet := s.alphabet[offset:] + s.alphabet[:offset]
	prefix := string(alphabet[0])
	partition := string(alphabet[1])
	alphabet = alphabet[2:]

	ret := []string{prefix}

	for i, num := range numbers {
		alphabetWithoutSeparator := alphabet[:len(alphabet)-1]
		ret = append(ret, toID(num, alphabetWithoutSeparator))

		if i < len(numbers)-1 {
			var separator string
			if partitioned && i == 0 {
				separator = partition
			} else {
				separator = string(alphabet[len(alphabet)-1])
			}

			ret = append(ret, separator)
			alphabet = shuffle(alphabet)
		}
	}

	id := strings.Join(ret, "")

	if s.minLength > len(id) {
		if !partitioned {
			numbers = append([]uint64{0}, numbers...)
			id, err = s.encodeNumbers(numbers, true)
			if err != nil {
				return "", err
			}
		}

		if s.minLength > len(id) {
			id = id[:1] + alphabet[:s.minLength-len(id)] + id[1:]
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

	alphabetChars := strings.Split(s.alphabet, "")

	for _, c := range strings.Split(id, "") {
		if !contains(alphabetChars, c) {
			return ret
		}
	}

	prefix := string(id[0])
	offset := strings.Index(s.alphabet, prefix)
	alphabet := s.alphabet[offset:] + s.alphabet[:offset]
	partition := string(alphabet[1])

	alphabet = alphabet[2:]
	id = id[1:]

	partitionIndex := strings.Index(id, partition)
	if partitionIndex > 0 && partitionIndex < len(id)-1 {
		id = id[partitionIndex+1:]
		alphabet = shuffle(alphabet)
	}

	for len(id) > 0 {
		separator := string(alphabet[len(alphabet)-1])
		chunks := strings.Split(id, separator)

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
				alphabet = shuffle(alphabet)
			}
		}

		id = strings.Join(chunks[1:], separator)
	}

	return ret
}

// MinValue returns the minimum uint64 value, which is 0
func MinValue() uint64 {
	return minUint64Value
}

// MaxValue returns the maximum uint64 value, which is 18446744073709551615
func MaxValue() uint64 {
	return maxUint64Value
}

func shuffle(alphabet string) string {
	chars := strings.Split(alphabet, "")

	for i, j := 0, len(chars)-1; j > 0; i, j = i+1, j-1 {
		r := (i*j + int(chars[i][0]) + int(chars[j][0])) % len(chars)
		chars[i], chars[r] = chars[r], chars[i]
	}

	return strings.Join(chars, "")
}

func toID(num uint64, alphabet string) string {
	id := []string{}
	chars := strings.Split(alphabet, "")

	result := num
	for {
		index := result % uint64(len(chars))

		id = append([]string{chars[index]}, id...)
		result = result / uint64(len(chars))

		if result == 0 {
			break
		}
	}

	return strings.Join(id, "")
}

func toNumber(id string, alphabet string) uint64 {
	chars := strings.Split(alphabet, "")
	result := uint64(0)

	for _, v := range id {
		result = result*uint64(len(chars)) + uint64(strings.Index(alphabet, string(v)))
	}

	return result
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

func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
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

func hasDigit(word string) bool {
	for _, r := range word {
		if r >= '0' && r <= '9' {
			return true
		}
	}

	return false
}
