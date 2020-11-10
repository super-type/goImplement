package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

// An example vendor B's backend, showing how easy it is to consume data from Supertype
func main() {
	// Gallagher's SupertypeID (will replace with a user's, not a vendor's)
	attribute := "masterBedroom"
	supertypeID := "gtgneBo6bnVpZC5lbGxpcHRpYy5jdXJ2ZS9wb2ludHgsQXpDTnpLVktQZWE3b2U4bXpKWnloTUgrRXlLZzBqbU1ibSthZFJCN2FmZmk="

	obs, err := goimplement.Consume(attribute, supertypeID, "MHcCAQEEIB0Co27xjk2xjBaZ4m5ebjscooulIAtxdjwVHJYAv4WDoAoGCCqGSM49AwEHoUQDQgAE4sPod+G1Nwfj11No5f2Qa2sUrTFTmoC4ppSfZrjg6YCPqb9ylaY+aBy1HeuM+8lhdB4CV2cvCV40yxBVy3kWag==", "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE4sPod+G1Nwfj11No5f2Qa2sUrTFTmoC4ppSfZrjg6YCPqb9ylaY+aBy1HeuM+8lhdB4CV2cvCV40yxBVy3kWag==", "zPti7f5IsUofeUjdBbH2+X3qKZqjfJRkOQ4RtpRNX87C9UFzXAnpqvGVXTKm9YwYA4vdnU7c8T3jNJ7TNL4hLFIvuhIsA7P/fZd6utZXcjcu9yqpUSbFB8SrWiuNkyDSPxAuiSWIuZ4=")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("obs: %v\n", *obs)
}
