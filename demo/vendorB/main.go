package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

// An example vendor B's backend, showing how easy it is to consume data from Supertype
func main() {
	attribute := "temperature"
	supertypeID := "user123"

	obs, err := goimplement.Consume(attribute, supertypeID, "72927570929357778628895897803930837315763393826720243708803951621364049710425", "045adbe11c882202334e20fee4466837037903eea4275832bccced0b8ce383179bd7eb84334222bcadddcfd930257fd9ba57f6b95e8ec0b7f77b667382a949a0c9")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("Consumed attribute %v from Supertype: %v\n", attribute, (*obs)[0].Plaintext)
}
