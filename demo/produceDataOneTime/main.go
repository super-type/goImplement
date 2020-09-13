package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

// An example vendor A's backend, showing how easy it is to produce data to Supertype
func main() {
	temperature := "200"
	attribute := "caloriesBurned"
	supertypeID := "user123"

	// this is the wrong sk for the given pk
	err := goimplement.Produce(temperature, attribute, supertypeID, "91032863405535923859153039920779555024328120618368196021389762057833204795281", "042e26494ef3eea3622d5bd7ba602e06476f4bde3f0418d2cb24c7d04fa7b6985105f2782cbc5a0bf7ac20911f7c7e307a9c77ce7ceef6e0cb45d4c996fc63a559")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("Produced data to Supertype")
}
