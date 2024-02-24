package progpjs

import (
	"github.com/progpjs/progpAPI/v2"
	"github.com/progpjs/progpAPI/v2/codegen"
)

type GenFCaller[T any] struct {
}

func (m *GenFCaller[T]) GetT() T {
	var v T
	return v
}

func GetFunctionCaller(functionTemplate any) any {
	sign := codegen.GetFunctionSignature(functionTemplate)

	// Get the final function.
	res := progpAPI.GetFunctionCaller(sign)
	if res != nil {
		return res
	}

	// The function doesn't exist?
	// Then will be added to the function to generate.
	//
	codegen.AddFunctionCallerToGenerate(functionTemplate)

	// Return the function used as template.
	// This function must contains a call to progpAPI.DynamicFunctionCaller
	// in order to be able to use dynamic mode.
	//
	return functionTemplate
}
