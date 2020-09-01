package main

import (
	"fmt"
	"log"

	"github.com/dnesting/unit/unitdef"
)

func main() {
	reg, err := unitdef.FromStandard()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", reg)
}
