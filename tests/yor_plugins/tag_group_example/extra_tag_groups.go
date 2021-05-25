package main

import "fmt"

var ExtraTagGroups = []interface{}{&OrgTagGroup{}}

func main() {
	fmt.Printf("We have %d tag groups here!", len(ExtraTagGroups))
}
