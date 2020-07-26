package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

func main() {
	// obs, err := goimplement.Consume("caloriesBurned", "test", "69088735009294891019519991270539942363091074413815484194427187441250221464654", "04d30747d3d66d288a9858eda09f8ffb1594b06e68d455975281be2ad25e3a31121298c49e6c5cd50a947d769ce241b0362ca6469e1a53d2bbf812f7673d36dea9")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	// fmt.Printf("Done! %v\n", obs)

	err := goimplement.Produce("hi guy, what's up", "caloriesBurned", "test", "91032863405535923859153039920779555024328120618368196021389762057833204795281", "04fcee14eaa83337315688ef3db63586ecceddfddf436b57444958806b566ec9d1e809577c44fcceeac9beb9d7be82f5b9be84d94f8bc54a82cf1a6f3c9ae37140")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Println("Done!")
	}
}
