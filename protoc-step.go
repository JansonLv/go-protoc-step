package protocstep

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"os"
	"path"
)

const defaultFileName = "input.bin"

type protocCodeGen struct {
	req     *pluginpb.CodeGeneratorRequest
	isWrite bool
}

// ReadFileWithJson 不建议使用序列化直接的json解析，因为会丢失Extension信息，不过可以作为参考和信息核对
func ReadFileWithJson(filename string) OptionFunc {
	if filename == "" {
		filename = "input.json"
	}
	return func(p *protocCodeGen) error {
		in, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		req := &pluginpb.CodeGeneratorRequest{}
		if err := protojson.Unmarshal(in, req); err != nil {
			return err
		}
		p.req = req
		return nil
	}
}

// ReadFileWithBin default proto-gen-step默认生成的文件名可不传
func ReadFileWithBin(filename string) OptionFunc {
	return func(p *protocCodeGen) error {
		in, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		req := &pluginpb.CodeGeneratorRequest{}
		if err := proto.Unmarshal(in, req); err != nil {
			return err
		}
		p.req = req
		return nil
	}
}

func WriteFile(isWrite bool) OptionFunc {
	return func(p *protocCodeGen) error {
		p.isWrite = isWrite
		return nil
	}
}

type OptionFunc func(gen *protocCodeGen) error

func ProtocStep(genCodeFunc func(gen *protogen.Plugin) error, opts ...OptionFunc) error {
	t := &protocCodeGen{isWrite: true}

	for _, opt := range opts {
		if err := opt(t); err != nil {
			return err
		}
	}
	if t.req == nil {
		if err := ReadFileWithBin(defaultFileName)(t); err != nil {
			return err
		}
	}

	gen, err := protogen.Options{}.New(t.req)
	if err != nil {
		return nil
	}
	err = genCodeFunc(gen)
	if err != nil {
		return err
	}
	response := gen.Response()
	if !t.isWrite {
		return nil
	}
	for _, file := range response.File {
		if err := write(file.GetName(), []byte(file.GetContent())); err != nil {
			return err
		}
	}
	return nil
}

func write(fileName string, bs []byte) error {
	absPath := path.Dir(fileName)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		if err := os.MkdirAll(absPath, 0777); err != nil {
			return err
		}
		// 再修改权限
		if err := os.Chmod(absPath, 0777); err != nil {
			return err
		}
	} else {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	if _, err = file.Write(bs); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
