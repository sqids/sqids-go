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

> 🚧 Please note that the following examples omit proper error handling:

Simple encode & decode:

```golang
s, _ := sqids.New()
id, _ := s.Encode([]uint64{1, 2, 3}) // "8QRLaD"
numbers := s.Decode(id) // [1, 2, 3]
```

Randomize IDs by providing a custom alphabet:

```golang
alphabet := "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE"
s, _ := sqids.NewCustom(sqids.Options{
    Alphabet: &alphabet,
})
id, _ := s.Encode([]uint64{1, 2, 3}) // "B5aMa3"
numbers := s.Decode(id) // [1, 2, 3]
```

Enforce a *minimum* length for IDs:

```golang
minLength := 10
s, _ := sqids.NewCustom(sqids.Options{
    MinLength: &minLength,
})
id, _ := s.Encode([]uint64{1, 2, 3}) // "75JT1cd0dL"
numbers := s.Decode(id) // [1, 2, 3]
```

Prevent specific words from appearing anywhere in the auto-generated IDs:

```golang
s, _ := sqids.NewCustom(sqids.Options{
    Blocklist: &[]string{"word1", "word2"},
})
id, _ := s.Encode([]uint64{1, 2, 3}) // "8QRLaD"
numbers := s.Decode(id) // [1, 2, 3]
```

## License

[MIT](LICENSE)
