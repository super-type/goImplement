package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

// An example vendor A's backend, showing how easy it is to produce data to Supertype
func main() {
	temperature := "153"
	attribute := "temperature"
	supertypeID := "user123"

	err := goimplement.Produce(temperature, attribute, supertypeID, "91032863405535923859153039920779555024328120618368196021389762057833204795281", "04fcee14eaa83337315688ef3db63586ecceddfddf436b57444958806b566ec9d1e809577c44fcceeac9beb9d7be82f5b9be84d94f8bc54a82cf1a6f3c9ae37140")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("Produced data to Supertype")
}
