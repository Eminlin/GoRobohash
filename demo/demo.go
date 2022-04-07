package main

import (
	"fmt"

	gorobohash "github.com/Eminlin/GoRobohash"
)

func main() {
	fmt.Println(gorobohash.NewResource("hello", &gorobohash.AssembleOptions{}).GenerateJPEG())
	// print: ./hello.png <nil>
	b64, err := gorobohash.NewResource("hello", &gorobohash.AssembleOptions{}).GenerateBase64()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b64))
}
