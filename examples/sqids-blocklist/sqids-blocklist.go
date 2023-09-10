package main

import (
	"fmt"

	"github.com/sqids/sqids-go"
)

func main() {
	s, _ := sqids.New(sqids.Options{
		Blocklist: []string{"86Rf07"},
	})
	id, _ := s.Encode([]uint64{1, 2, 3}) // "se8ojk"
	numbers := s.Decode(id)              // [1, 2, 3]

	fmt.Println(id, numbers)
}
