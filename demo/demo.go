package main

import (
	gorobohash "GoRobohash"
	"fmt"
)

func main() {
	fmt.Println(gorobohash.NewResource("hello", &gorobohash.AssembleOptions{}).Assemble())
	// print: ./hello.png <nil>
}
