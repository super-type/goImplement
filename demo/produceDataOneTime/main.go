package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

// An example vendor A's backend, showing how easy it is to produce data to Supertype
func main() {
	// Gallagher's SupertypeID (will replace with a user's, not a vendor's)
	temperature := "200"
	attribute := "temperature"
	supertypeID := "gtgneBo6bnVpZC5lbGxpcHRpYy5jdXJ2ZS9wb2ludHgsQWh0S0V1cWlJTlJrTnZyMWRBTHdWNWxNRzNoVVJTUHduNGhUZmh2Qm44c0E="

	err := goimplement.Produce(temperature, attribute, supertypeID, "MHcCAQEEIKn/I4RaVf7/p3QbYqwH0nQJsjRqKwn/7YUJ/eljNMwroAoGCCqGSM49AwEHoUQDQgAEDywt/8GOiJHxa7yY1l/fYj0Y3p6ITIhh5LqlwMtGjd8Wiy7bx4eY3FsoRKtb1CRlYhGOFsb8Se7Ya2VcqJfecA==", "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEDywt/8GOiJHxa7yY1l/fYj0Y3p6ITIhh5LqlwMtGjd8Wiy7bx4eY3FsoRKtb1CRlYhGOFsb8Se7Ya2VcqJfecA==")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("Produced data to Supertype")
}
