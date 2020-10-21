package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

func main() {
	// Gallagher's SupertypeID (will replace with a user's, not a vendor's)
	attribute := "temperature"
	supertypeID := "gtgneBo6bnVpZC5lbGxpcHRpYy5jdXJ2ZS9wb2ludHgsQWh0S0V1cWlJTlJrTnZyMWRBTHdWNWxNRzNoVVJTUHduNGhUZmh2Qm44c0E="

	err := goimplement.ConsumeWS(attribute, supertypeID, "MHcCAQEEIB0Co27xjk2xjBaZ4m5ebjscooulIAtxdjwVHJYAv4WDoAoGCCqGSM49AwEHoUQDQgAE4sPod+G1Nwfj11No5f2Qa2sUrTFTmoC4ppSfZrjg6YCPqb9ylaY+aBy1HeuM+8lhdB4CV2cvCV40yxBVy3kWag==", "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE4sPod+G1Nwfj11No5f2Qa2sUrTFTmoC4ppSfZrjg6YCPqb9ylaY+aBy1HeuM+8lhdB4CV2cvCV40yxBVy3kWag==")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
