package sqids

//go:generate go run github.com/campoy/embedmd/v2@v2.0.0 -w README.md

import (
	"errors"
	"strings"
)

const (
	defaultAlphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	minAlphabetLength = 3
)

var defaultBlocklist []string = newDefaultBlocklist()

// Alphabet validation errors
var (
	errAlphabetMultibyte       = errors.New("alphabet must not contain any multibyte characters")
	errAlphabetTooShort        = errors.New("alphabet length must be at least 3")
	errAlphabetNotUniqueChars  = errors.New("alphabet must contain unique characters")
	errMaxRegenerationAttempts = errors.New("reached max attempts to re-generate the id")
)

// Options for a custom instance of Sqids
type Options struct {
	Alphabet  string
	MinLength uint8
	Blocklist []string
}

// Sqids lets you generate unique IDs from numbers
type Sqids struct {
	alphabet  string
	minLength uint8
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

	o.Blocklist = filterBlocklist(o.Alphabet, o.Blocklist)

	return o, nil
}

// Encode a slice of uint64 values into an ID string
func (s *Sqids) Encode(numbers []uint64) (string, error) {
	// if no numbers passed, return an empty string
	if len(numbers) == 0 {
		return "", nil
	}

	return s.encodeNumbers(numbers, 0)
}

func (s *Sqids) encodeNumbers(numbers []uint64, increment int) (string, error) {
	if increment > len(s.alphabet) {
		return "", errMaxRegenerationAttempts
	}

	var (
		err      error
		offset   = calculateOffset(s.alphabet, numbers, increment)
		alphabet = alphabetOffset(s.alphabet, offset)
		prefix   = alphabet[0]
		ret      = []rune{prefix}
	)

	alphabet = reverseRunes(alphabet)

	for i, num := range numbers {
		ret = append(ret, []rune(toID(num, string(alphabet[1:])))...)

		if i < len(numbers)-1 {
			ret = append(ret, alphabet[0])
			alphabet = []rune(shuffle(string(alphabet)))
		}
	}

	id := string(ret)

	if int(s.minLength) > len(id) {
		id += string(alphabet[0])

		for int(s.minLength)-len(id) > 0 {
			alphabet = []rune(shuffle(string(alphabet)))
			id += string(alphabet[:min(int(s.minLength)-len(id), len(alphabet))])
		}
	}

	if s.isBlockedID(id) {
		id, err = s.encodeNumbers(numbers, increment+1)
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
	alphabet = reverseRunes(alphabet)

	rid = rid[1:]

	for len(rid) > 0 {
		separator := alphabet[0]

		chunks := splitChunks(rid, separator)
		if len(chunks) > 0 {
			if len(chunks[0]) == 0 {
				return ret
			}

			ret = append(ret, toNumber(chunks[0], alphabet[1:]))

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
	var chunks [][]rune
	chunk := []rune{}

	for _, r := range runes {
		if r == separator {
			chunks = append(chunks, chunk)
			chunk = []rune{}
		} else {
			chunk = append(chunk, r)
		}
	}

	chunks = append(chunks, chunk)
	return chunks
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

func calculateOffset(alphabet string, numbers []uint64, increment int) int {
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

	offset = offset % len(runes)
	return (offset + increment) % len(runes)
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

func reverseRunes(runes []rune) []rune {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
