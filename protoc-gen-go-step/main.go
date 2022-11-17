package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	// api常用的，默认包含解析
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

func main() {
	var binFileName = flag.String("binfile", "input.bin", "字节流传输文件")
	var jsonFileName = flag.String("jsonfile", "input.json", "序列化协议")

	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}

	// 将序列化的协议解析出flag
	gen, err := protogen.Options{ParamFunc: flag.Set}.New(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		return
	}
	// 支持optional，否则os.Stdout.Write会显示相关问题
	//a proto3 file that contains optional fields,
	//but code generator protoc-gen-data hasn't been updated to support optional fields in proto3.
	//Please ask the owner of this code generator to support proto3 optional.--data_out:
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	resp := gen.Response()

	marshal, _ := protojson.Marshal(req)
	resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    jsonFileName,
		Content: proto.String(string(marshal)),
	})
	resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    binFileName,
		Content: proto.String(string(in)),
	})

	out, err := proto.Marshal(resp)
	if err != nil {
		return
	}
	if _, err := os.Stdout.Write(out); err != nil {
		return
	}
}
