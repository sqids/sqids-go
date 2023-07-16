# [Sqids Go](https://sqids.org/go)

[![Github Actions](https://img.shields.io/github/actions/workflow/status/sqids/sqids-go/tests.yml)](https://github.com/sqids/sqids-go/actions)

Sqids (pronounced "squids") is a small library that lets you generate YouTube-looking IDs from numbers. It's good for link shortening, fast & URL-safe ID generation and decoding back into numbers for quicker database lookups.

## Getting started

Use go get.

```bash
go get github.com/sqids/sqids-go
```

Then import the package into your own code.

```golang
import "github.com/sqids/sqids-go"
```

## Examples

Simple encode & decode:

```golang
s, err := sqids.New()
id := s.Encode([]uint64{1, 2, 3}) // "8QRLaD"
numbers := s.Decode(id) // [1, 2, 3]
```

Randomize IDs by providing a custom alphabet:

```golang
s, err := sqids.NewCustom("FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE", 0)
id := s.Encode([]uint64{1, 2, 3}) // "B5aMa3"
numbers := s.Decode(id) // [1, 2, 3]
```

Enforce a *minimum* length for IDs:

```golang
s, err := sqids.NewCustom("", 10)
id := s.Encode([]uint64{1, 2, 3}) // "75JT1cd0dL"
numbers := s.Decode(id) // [1, 2, 3]
```

Prevent specific words from appearing anywhere in the auto-generated IDs:

```golang
s, err := sqids.NewSqids(sqids.Options{
    Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
    MinLength: 0,
    Blocklist: &[]string{"word1", "word2"}
})
id := s.Encode([]uint64{1, 2, 3}) // "8QRLaD"
numbers := s.Decode(id) // [1, 2, 3]
```

## License

[MIT](LICENSE)
