package main

import (
	"fmt"

	"github.com/sqids/sqids-go"
)

func main() {
	s, _ := sqids.New(sqids.Options{
		Alphabet: "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE",
	})
	id, _ := s.Encode([]uint64{1, 2, 3}) // "B4aajs"
	numbers := s.Decode(id)              // [1, 2, 3]

	fmt.Println(id, numbers)
}
