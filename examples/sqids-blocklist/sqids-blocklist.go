package main

import (
	"fmt"

	"github.com/sqids/sqids-go"
)

func main() {
	s, _ := sqids.New(sqids.Options{
		Blocklist: []string{"word1", "word2"},
	})
	id, _ := s.Encode([]uint64{1, 2, 3}) // "8QRLaD"
	numbers := s.Decode(id)              // [1, 2, 3]

	fmt.Println(id, numbers)
}
