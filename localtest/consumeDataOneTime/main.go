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

	obs, err := goimplement.Consume(attribute, supertypeID, "MHcCAQEEIJ8nXdJhT4jHQVmGleWz5lSFQXPGiOCakWT2SY1zovPHoAoGCCqGSM49AwEHoUQDQgAEiSObPMZCzxrFYwn7JC3qdIsVKhE/KnK/n/V8/3z6DFg+evII+2IRwkATwOfukeTFyMM8hUGjCqvvT4deq+PmTw==", "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEiSObPMZCzxrFYwn7JC3qdIsVKhE/KnK/n/V8/3z6DFg+evII+2IRwkATwOfukeTFyMM8hUGjCqvvT4deq+PmTw==", "zPti7f5IsUofeUjdBbH2+X3qKZqjfJRkOQ4RtpRNX87C9UFzXAnpqvGVXTKm9YwYA4vdnU7c8T3jNJ7TNL4hLFIvuhIsA7P/fZd6utZXcjcu9yqpUSbFB8SrWiuNkyDSPxAuiSWIuZ4=")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("obs: %v\n", *obs)
}
