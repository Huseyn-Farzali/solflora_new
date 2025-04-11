package main

import "fmt"

type Guarantee struct {
	PrincipalName   string
	GuaranteeAmount float64
}

func main() {
	var guarantee = new(Guarantee)
	guarantee.PrincipalName = "John"
	fmt.Println(*guarantee)
}
