package main

import (
	"fmt"

	"github.com/sqids/sqids-go"
)

func main() {
	s, _ := sqids.New()
	id, _ := s.Encode([]uint64{1, 2, 3}) // "86Rf07"
	numbers := s.Decode(id)              // [1, 2, 3]

	fmt.Println(id, numbers)
}
