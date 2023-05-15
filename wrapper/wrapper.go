package wrapper

import (
	"log"
	"unsafe"

	"github.com/lingkaix/demo-ns/protos/mydata"
)

var handler (func(ctx mydata.Input) mydata.Output)

func Handle(f func(ctx mydata.Input) mydata.Output) {
	handler = f
}

//export handle_http
func _handle(dataPtr, size int32) int64 {
	reqData := ptr2bytes(dataPtr, size)
	req := &mydata.Input{}
	req.UnmarshalVT(reqData)

	resp := handler(*req)
	// ! is b in a linear memory allocation?
	b, err := resp.MarshalVT()
	if err != nil {
		log.Println(err)
	}
	return bytes2ptr(b)
}

func ptr2bytes(dataPtr, size int32) []byte {
	p := unsafe.Pointer(uintptr(dataPtr))
	s := uint(size)
	b := ((*[1 << 30]byte)(p))[:s:s]

	return b
}

func bytes2ptr(data []byte) int64 {
	return ((int64(uintptr(unsafe.Pointer(&data[0]))) << int64(32)) | int64(len(data)))
}
