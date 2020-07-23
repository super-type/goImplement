package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

func main() {
	err := goimplement.Produce("hi guy, what's up", "caloriesBurned", "test", "66781207567065100428315162405958259237703165062370629758395207375855887994745", "045c07f767fdb7cb9c9646c2c8d22ff8079a61c011e8c3bfd13b6e172020ec7be2bbf1a33dbd656ee30ab91ca6318d2a00b0649f8dc009d6ff83a103e9265f2e86")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Println("Done!")
	}
	// err := goimplement.Consume("caloriesBurned", "test", "69088735009294891019519991270539942363091074413815484194427187441250221464654", "04d30747d3d66d288a9858eda09f8ffb1594b06e68d455975281be2ad25e3a31121298c49e6c5cd50a947d769ce241b0362ca6469e1a53d2bbf812f7673d36dea9")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	// fmt.Println("Done!")
}
