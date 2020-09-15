package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

// An example vendor A's backend, showing how easy it is to produce data to Supertype
func main() {
	temperature := "200"
	attribute := "temperature"
	supertypeID := "user123"

	err := goimplement.Produce(temperature, attribute, supertypeID, "40564280015546133127186018201877113429491387889238716460773792707727598260818", "04eca9461980bc478f86c56ca316a959e283c4be77da6a63aca2b26a349e103b24b0595affd6ce401d323f6fc61c003cd580f609b7f9ec58901259ce1190aecc0a")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("Produced data to Supertype")
}
