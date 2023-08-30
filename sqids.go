package sqids

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	// defaultAlphabet -
	defaultAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// defaultMinLength -
	defaultMinLength int = 0

	// defaultBlocklist -
	defaultBlocklist []string = newDefaultBlocklist()
)

// Options -
type Options struct {
	Alphabet  *string
	MinLength *int
	Blocklist *[]string
}

// Sqids -
type Sqids struct {
	alphabet  string
	minLength int
	blocklist []string
}

// New -
func New() (*Sqids, error) {
	return NewCustom(Options{
		Alphabet:  &defaultAlphabet,
		MinLength: &defaultMinLength,
		Blocklist: &defaultBlocklist,
	})
}

// NewCustom -
func NewCustom(options Options) (*Sqids, error) {
	alphabet := options.Alphabet
	if alphabet == nil {
		alphabet = &defaultAlphabet
	}

	minLength := options.MinLength
	if minLength == nil {
		minLength = &defaultMinLength
	}

	blocklist := options.Blocklist
	if blocklist == nil {
		blocklist = &defaultBlocklist
	}

	// check the length of the alphabet
	if len(*alphabet) < 5 {
		return nil, errors.New("alphabet length must be at least 5")
	}

	// check that the alphabet has only unique characters
	if !hasUniqueChars(*alphabet) {
		return nil, errors.New("alphabet must contain unique characters")
	}

	// test min length (type [might be lang-specific] + min length + max length)
	if *minLength < int(MinValue()) || *minLength > len(*alphabet) {
		return nil, fmt.Errorf("minimum length has to be between %d and %d", MinValue(), len(*alphabet))
	}

	// clean up blocklist:
	// 1. all blocklist words should be lowercase
	// 2. no words less than 3 chars
	// 3. if some words contain chars that are not in the alphabet, remove those
	filteredBlocklist := []string{}
	alphabetChars := strings.Split(strings.ToLower(*alphabet), "")
	for _, word := range *blocklist {
		if len(word) >= 3 {
			wordLowercased := strings.ToLower(word)
			wordChars := strings.Split(wordLowercased, "")
			intersection := intersection(wordChars, alphabetChars)
			if len(intersection) == len(wordChars) {
				filteredBlocklist = append(filteredBlocklist, strings.ToLower(wordLowercased))
			}
		}
	}

	return &Sqids{
		alphabet:  shuffle(*alphabet),
		minLength: *minLength,
		blocklist: filteredBlocklist,
	}, nil
}

// Encode -
func (s *Sqids) Encode(numbers []uint64) (string, error) {
	// if no numbers passed, return an empty string
	if len(numbers) == 0 {
		return "", nil
	}

	inRangeNumbers := []uint64{}
	for _, n := range numbers {
		if n >= MinValue() && n <= MaxValue() {
			inRangeNumbers = append(inRangeNumbers, n)
		}
	}

	if len(inRangeNumbers) != len(numbers) {
		return "", fmt.Errorf("encoding supports numbers between %d and %d", MinValue(), MaxValue())
	}

	return s.encodeNumbers(inRangeNumbers, false)
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
			if numbers[0]+1 > MaxValue() {
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

// Decode -
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

// MinValue -
func MinValue() uint64 {
	return 0
}

// MaxValue -
func MaxValue() uint64 {
	return math.MaxUint64
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

func intersection(slice1, slice2 []string) []string {
	intersect := []string{}
	set := make(map[string]bool)

	for _, s := range slice2 {
		set[s] = true
	}

	for _, s := range slice1 {
		if set[s] {
			intersect = append(intersect, s)
		}
	}

	return intersect
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
