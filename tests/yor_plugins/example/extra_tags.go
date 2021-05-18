package main

import "fmt"

var ExtraTags = []interface{}{&GitOwnerTag{}, &FooTag{}}

func main() {
	fmt.Printf("We have %d tags here!", len(ExtraTags))
}
