package main

import (
	"github.com/lingkaix/demo-ns/protos/mydata"
	"github.com/lingkaix/demo-ns/wrapper"
)

func init() {
	wrapper.Handle(Recv)
}
func main() {}

func Recv(req mydata.Input) mydata.Output {

	sum := int32(0)
	for _, v := range req.GetIntValues() {
		sum += v
	}
	str := ""
	for _, v := range req.GetStringValues() {
		str = str + v
	}
	return mydata.Output{IntResult: sum, StringResult: str}
}
