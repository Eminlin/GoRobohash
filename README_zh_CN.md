# GoRobohash
Robohash 项目 Go 版本, 原版： [e1ven/Robohash](https://github.com/e1ven/Robohash)
# Readme
[[English](./README.md)] [[简体中文](./README_zh_CN.md)]

# 主要目录文件
|-- `material/...` 生成图片所需要的素材文件： [License](https://github.com/e1ven/Robohash#license)  
|-- `robohash.go`  代码来源: [robohash.py](https://github.com/e1ven/Robohash/blob/master/robohash/robohash.py)

## 使用方法

```go
go get -u github.com/Eminlin/GoRobohash
```

## 示例

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

然后你就在运行根目录得到一张和 [robohash.org](https://robohash.org/) 一样的照片。


更多使用方式可以在 `demo.go` 里面查看。

## 配置 

```go
type AssembleOptions struct {
	RoboSet, Colors, BgSet string //可选
	OutputPath             string //可选 默认当前路径
	X                      int    //可选 默认 300
	Y                      int    //可选 默认 300
}
```

## 提示

原始项目携带很多 png 静态资源文件，如果需要 `go build` 在其它机器运行，要保证机器里有存在对应的资源文件，建议直接在对应机器执行 `go get -u github.com/Eminlin/GoRobohash`  

目前不支持 `bmp` 格式的图片。  

建议生成 `png` 格式的图片。