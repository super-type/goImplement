package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

func main() {
	attribute := "temperature"
	supertypeID := "user123"

	err := goimplement.ConsumeWS(attribute, supertypeID, "8575e5f7ca6ed6345ccdfc2d595649c571c6a4ac9a847b049b1ea67eb5e53acf", "046a9d202dcad0c86aa88fb72472fa2b2f180893bc4802a1fe9ce6ae80b4c1cef9a89d533c047e391f778759460830db149f9944c32126e4827c48819986238315")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
