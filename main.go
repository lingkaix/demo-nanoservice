package main

import (
	"fmt"
	"os"

	"github.com/lingkaix/demo-ns/protos/mydata"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func main() {
	wasmBytes, _ := os.ReadFile("wasm.wasm")
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	fmt.Println("Compiling module...")
	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		fmt.Println("Failed to compile module:", err)
	}
	wasiEnv, _ := wasmer.NewWasiStateBuilder("wasi-program").
		// Choose according to your actual situation
		// Argument("--foo").
		// Environment("ABC", "DEF").
		// MapDirectory("./", ".").
		Finalize()
	importObject, err := wasiEnv.GenerateImportObject(store, module)
	check(err)
	// Create an empty import object.
	// importObject := wasmer.NewImportObject()

	fmt.Println("Instantiating module...")
	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		panic(fmt.Sprintln("Failed to instantiate the module:", err))
	}

	memory, err := instance.Exports.GetMemory("memory")
	if err != nil {
		panic(fmt.Sprintln("Failed to get the `memory` memory:", err))
	}
	memory.Grow(20)
	memArr := memory.Data()

	malloc, err := instance.Exports.GetFunction("malloc")
	if err != nil {
		panic(fmt.Sprintln("Failed to retrieve the `malloc` function:", err))
	}
	free, err := instance.Exports.GetFunction("free")
	if err != nil {
		panic(fmt.Sprintln("Failed to retrieve the `free` function:", err))
	}
	handler, err := instance.Exports.GetFunction("handle_http")
	if err != nil {
		panic(fmt.Sprintln("Failed to retrieve the `handler` function:", err))
	}
	req := request()

	size := int32(len(req))
	input, err := malloc(size)
	defer free(input)
	if err != nil {
		fmt.Println(err)
	}
	ptr := input.(int32)
	copy(memArr[ptr:], req)

	output, err := handler(ptr, size)
	if err != nil {
		fmt.Println(err)
	}
	ptrSize := output.(int64)
	respPtr, respSize := splitint64(ptrSize)
	respBytes := memArr[respPtr : respPtr+respSize]
	resp := &mydata.Output{}
	err = resp.UnmarshalVT(respBytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.GetStringResult(), resp.GetIntResult())
}

func request() []byte {
	intValues := []int32{1, 2, 3}
	stringValues := []string{"Hello", "World"}

	// Create the Input message
	inputMessage := &mydata.Input{
		IntValues:    intValues,
		StringValues: stringValues,
	}

	// Serialize the Input message
	data, err := inputMessage.MarshalVT()
	if err != nil {
		return nil
		// fmt.Fatal("Failed to marshal input message:", err)
	}
	return data
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func splitint64(n int64) (int32, int32) {
	lower := int32(n)
	upper := int32(n >> 32)
	return upper, lower
}
