package main

import (
	"fmt"
	"temporal/greeting"
)

func main() {
	greeter := greeting.Greet("Nolan")
	fmt.Println(greeter)
}
