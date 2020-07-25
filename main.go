package main

import (
	"fmt"

	goimplement "github.com/super-type/supertype/goImplement"
)

func main() {
	// err := goimplement.Produce("hi guy, what's up", "caloriesBurned", "test", "90108926111724530182196785119666121890754471775626444759124217906070508658767", "04528eec550ffcb5f1451f8d9d995444390f87a87d4a6c3d0586e92cbc159c64ac04c4785d45b06f34707355fd0433d67c9a7e64867d71358c910957c67fcf1a0d")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// } else {
	// 	fmt.Println("Done!")
	// }
	// obs, err := goimplement.Consume("caloriesBurned", "test", "69088735009294891019519991270539942363091074413815484194427187441250221464654", "04d30747d3d66d288a9858eda09f8ffb1594b06e68d455975281be2ad25e3a31121298c49e6c5cd50a947d769ce241b0362ca6469e1a53d2bbf812f7673d36dea9")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	// fmt.Printf("Done! %v\n", obs)

	err := goimplement.Produce("hi guy, what's up", "caloriesBurned", "test", "69088735009294891019519991270539942363091074413815484194427187441250221464654", "04d30747d3d66d288a9858eda09f8ffb1594b06e68d455975281be2ad25e3a31121298c49e6c5cd50a947d769ce241b0362ca6469e1a53d2bbf812f7673d36dea9")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Println("Done!")
	}
}
