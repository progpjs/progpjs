package progpjs

import (
	"github.com/progpjs/progpAPI/v2"
	"github.com/progpjs/progpAPI/v2/codegen"
	"reflect"
)

type GenFCaller[T any] struct {
}

func (m *GenFCaller[T]) GetT() T {
	var v T
	return v
}

func GetFunctionCaller(defaultImpl any) any {
	reflectI := reflect.TypeOf(defaultImpl)
	myMethod, isFound := reflectI.MethodByName("Call")

	if !isFound {
		panic("Interface must have a methode named 'Call'")
	}

	myMethodT := myMethod.Type

	// Must always be added to the code generator.
	codegen.AddFunctionCallerToGenerate(myMethodT)

	sign := codegen.GetFunctionSignatureWithoutReturn(myMethodT)

	// Get the final function.
	res := progpAPI.GetFunctionCaller(sign)
	//
	if res == nil {
		return defaultImpl
	}

	return res
}
