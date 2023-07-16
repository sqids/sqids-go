package sqids

const (
	// DefaultAlphabet is the default alphabet
	DefaultAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// DefaultMinLength is the minimum length of generated IDs
	DefaultMinLength int = 0
)

// Sqids -
type Sqids struct {
	// Alphabet is the alphabet used to generate new IDs
	Alphabet string

	// MinLength is the minimum length of generated IDs
	MinLength int

	// Blocklist contains a list of words that should not appear in generated IDs
	Blocklist []string
}

// New creates a new Sqids with default parameters
func New() (*Sqids, error) {
	return &Sqids{Alphabet: DefaultAlphabet, MinLength: DefaultMinLength, Blocklist: Blocklist()}, nil
}

// Encode -
func (s *Sqids) Encode(numbers []uint64) (string, error) {
	return "", nil
}

// Decode -
func (s *Sqids) Decode(id string) ([]uint64, error) {
	return []uint64{}, nil
}

// MinValue -
func (s *Sqids) MinValue() uint64 {
	return 0
}

// MaxValue -
func (s *Sqids) MaxValue() uint64 {
	return uint64(^uint(0))
}
