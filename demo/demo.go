package main

import (
	"fmt"

	gorobohash "github.com/Eminlin/GoRobohash"
)

func main() {
	r, err := gorobohash.NewResource("hello", &gorobohash.AssembleOptions{
		X: 500,
		Y: 500,
	}).GeneratePNG()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
	// print: ./hello.png <nil>

	r, err = gorobohash.NewResource("hello", &gorobohash.AssembleOptions{}).GenerateJPEG()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
	// print: ./hello.jpeg <nil>

	b64, err := gorobohash.NewResource("hello", &gorobohash.AssembleOptions{}).GenerateBase64()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b64))
	// get base64 image:iVBORw0KGgoAAAA...48AAAAASUVORK5CYII= <nil>
}
