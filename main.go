package main

import (
	"academy/patterns"
	"fmt"
)

func main() {
	p, _ := patterns.Compile("abc*")
	fmt.Println(p.Match("abc"))
}
