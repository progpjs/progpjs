package scriptTransformer

func SetJavascriptModuleResolver(handler JavascriptModuleResolverF) {
	gJavascriptModuleResolver = handler
}
