# GoRobohash
Robohash Project Go Version, origin from [e1ven/Robohash](https://github.com/e1ven/Robohash)

# Readme
[[English](./README.md)] [[简体中文](./README_zh_CN.md)]

# Main File Tree
|-- `material/...` Materials needed to assemble pictures: [License](https://github.com/e1ven/Robohash#license)  
|-- `robohash.go`  Origin Code: [robohash.py](https://github.com/e1ven/Robohash/blob/master/robohash/robohash.py)

## How to use

```go
go get -u github.com/Eminlin/GoRobohash
```

## Example

```go
import (
	gorobohash "github.com/Eminlin/GoRobohash"
)

func main() {
	r, err := gorobohash.NewResource("hello", &gorobohash.AssembleOptions{}).GeneratePNG()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
    // print: ./hello.png <nil>
}
```

Then you can get a picture with the same style as [robohash.org](https://robohash.org/)

## Config

```go
type AssembleOptions struct {
    RoboSet, Colors, BgSet string //optional
    OutputPath             string //optional default current path
    X                      int    //optional default 300
    Y                      int    //optional default 300
}
```

More usage can be found in `demo.go`.

## Tips
Original project has many static resources, if you need `go build` on other machines, you need to ensure the resources exist on the machine. I recommend you to directly execute `go get -u github.com/Eminlin/GoRobohash` on the machine.  

Not support `bmp` format picture yet.

Recommend generate `png` format picture.