package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

// An example vendor B's backend, showing how easy it is to consume data from Supertype
func main() {
	attribute := "caloriesBurned"
	supertypeID := "user123"

	// TODO when making a new vendor account, add the skB, pkB
	obs, err := goimplement.Consume(attribute, supertypeID, "", "")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("Consumed observation %v: %v\n", attribute, obs)
}
