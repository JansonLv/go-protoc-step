# protoc-gen-go-step

protoc插件单步调试工具包

### step1. 安装

```shell
# 安装protoc-gen插件
go install github.com/JansonLv/go-protoc-step/protoc-gen-go-step@latest

# 安装调试包
go get -u github.com/JansonLv/go-protoc-step
```

### step2. 生成input文件

```shell
protoc --proto_path=./api --proto_path=./third_party \
       --go-step_out=paths=source_relative:./api  
       api/xxxx.proto
```

该命令会生成两个文件，一个是json文件，一个是bin文件。

通过设置binfile和jsonfile参数设置数据json和bin文件的文件名:

```
--go-step_out=paths=source_relative,jsonfile=xxx.json:./api \
```

bin其实就是protoc通过os.stdin输入到插件中的字节流数据。

json文件时这些字节流通过protojson序列化生成的文件，但是可能会丢失一些插件定义的Extension信息

后面单步调试时请尽量使用bin文件进行调试，json文件可以做一些非option下的数据调试或者阅读。

### step3. 单步调试

利用go test建立测试函数，将step2的bin或者json文件拷贝到当前目录下，调用ProtocStep函数，接下来就看你的啦

```go
package main

import (
	"github.com/JansonLv/go-protoc-step"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"testing"
)

func Test_generateFile(t *testing.T) {
	err := protocstep.ProtocStep(func(gen *protogen.Plugin) error {
		for _, file := range gen.Files {
			if file.Generate {
				generateFile(gen, file)
			}
		}
		return nil
	})
	assert.NoError(t, err)
}
```

ProtocStep默认使用input.bin文件，如需使用json文件或者其他文件名，可以使用其option模式

```go
// 使用xxx.bin作为输入文件
err := protocstep.ProtocStep(f, protocstep.ReadFileWithBin("xxx.bin"))
// 使用xxx.json输入文件
err = protocstep.ProtocStep(f, protocstep.ReadFileWithJson("xxx.json"))
// 运行完不生成go文件
err = protocstep.ProtocStep(f, protocstep.WriteFile(false))
// 组合
err = protocstep.ProtocStep(f, protocstep.ReadFileWithBin("xxx.bin"), protocstep.WriteFile(false))
```
