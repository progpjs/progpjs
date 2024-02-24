package progpjs

import "github.com/progpjs/progpAPI/v2/codegen"

func GetFunctionCaller[K any](functionTemplate K) K {
	res := codegen.GetFunctionCaller(functionTemplate)

	if res == nil {
		return functionTemplate
	}

	return res.(K)
}
