package main

import (
	"fmt"

	"github.com/sqids/sqids-go"
)

func main() {
	s, _ := sqids.New(sqids.Options{
		MinLength: 10,
	})
	id, _ := s.Encode([]uint64{1, 2, 3}) // "75JT1cd0dL"
	numbers := s.Decode(id)              // [1, 2, 3]

	fmt.Println(id, numbers)
}
