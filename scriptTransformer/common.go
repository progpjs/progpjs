package scriptTransformer

type TransformedScript struct {
	CompiledScriptPath    string
	OutputDir             string
	CompiledScriptContent string
	SourceMapFileContent  string
	SourceMapScriptPath   string
}
