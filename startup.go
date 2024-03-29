/*
 * (C) Copyright 2024 Johan Michel PIQUET, France (https://johanpiquet.fr/).
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package progpjs

import (
	"github.com/go-sourcemap/sourcemap"
	"github.com/progpjs/progpAPI/v2"
	"github.com/progpjs/progpAPI/v2/codegen"
	"github.com/progpjs/progpjs/v2/scriptTransformer"
	"os"
	"path"
	"plugin"
	"runtime"
)

//region Script engine resolver

//region Go plugin loader

var gGoPluginAreLoaded bool

func logGoPluginError(err error, pluginPath string) {
	if err == nil {
		return
	}

	//runtime.Breakpoint()

	println("Can't load plugin", pluginPath)
	println("Error:", err.Error())
	os.Exit(1)
}

func loadGoPlugin(pluginPath string) {
	// Warning: if send error when using debugger.
	// You must build your Go code with options 		-gcflags='all=-N -l'
	// which disable inlining and optimisations.
	//
	// Same for the plugin:
	// go build -buildmode=plugin -gcflags='all=-N -l' -o ./plugins/progpV8.so ./progpgo.scriptEngine.progpV8/asInstaller/installer.go

	_, err := plugin.Open(pluginPath)
	logGoPluginError(err, pluginPath)

	// About plugins:
	// - Must be in "main" package.
	// - Once loaded, the "func init() { ... }" is called.
}

func loadEngineFromPlugins(pluginsDir string) {
	if gGoPluginAreLoaded {
		return
	}
	gGoPluginAreLoaded = true

	loadGoPlugin(path.Join(pluginsDir, "progpV8.so"))
}

//endregion

func resolveScriptEngine(engineName string, pluginsDir string) progpAPI.ScriptEngine {
	if engineName == "" {
		engineName = "progpV8"
	}

	if gScriptEngine != nil {
		return gScriptEngine
	}

	gScriptEngine = progpAPI.UseScriptEngine(engineName)
	if gScriptEngine != nil {
		return gScriptEngine
	}

	loadEngineFromPlugins(pluginsDir)
	gScriptEngine = progpAPI.UseScriptEngine(engineName)
	if gScriptEngine != nil {
		return gScriptEngine
	}

	println("No script engine found !!!")
	os.Exit(1)

	return nil
}

var gScriptEngine progpAPI.ScriptEngine

//endregion

//region Script resources resolver

func resolveMissingJavascriptModule(resourceName string) (content string, loader uint16, isFound bool) {
	provider := gJavascriptModuleProviders[resourceName]
	if provider == nil {
		return "", 0, false
	}

	isFound = true
	var tLoader JsResourceLoader
	content, tLoader = provider(resourceName)
	loader = uint16(tLoader)

	return
}

type JsResourceLoader uint16

// This value must mimic esbuild "api.Loader" values.
const (
	JsLoaderNone JsResourceLoader = iota
	JsLoaderBase64
	JsLoaderBinary
	JsLoaderCopy
	JsLoaderCSS
	JsLoaderDataURL
	JsLoaderDefault
	JsLoaderEmpty
	JsLoaderFile
	JsLoaderGlobalCSS
	JsLoaderJS
	JsLoaderJSON
	JsLoaderJSX
	JsLoaderLocalCSS
	JsLoaderText
	JsLoaderTS
	JsLoaderTSX
)

//endregion

//region Config items

type EngineOptions struct {
	MustDebug               bool
	ScriptEngineName        string
	PluginsDir              string
	ProgpV8EngineProjectDir string

	OnScriptCompilationError  progpAPI.OnScriptCompilationErrorF
	OnRuntimeError            progpAPI.RuntimeErrorHandlerF
	OnScriptTerminated        progpAPI.ScriptTerminatedHandlerF
	OnCheckingAllowedFunction progpAPI.CheckAllowedFunctionsF

	// ScriptAutomaticHeader allows to insert an head at the start of each script.
	// If blank, then use the default value which is : import '@progp/core'
	// If you don"t want an header the use a string with only a space or "//".
	ScriptAutomaticHeader string
}

var gJavascriptModuleProviders = make(map[string]JavascriptModuleProviderF)

type JavascriptModuleProviderF func(resourcePath string) (content string, loader JsResourceLoader)

//endregion

var gIsBootstrapped = false
var gEngineOptions *EngineOptions
var gDefaultScriptEngine progpAPI.ScriptEngine
var gEnableDebugger bool

func GetScriptEngine() progpAPI.ScriptEngine {
	return gDefaultScriptEngine
}

func compileScript(scriptPath string) (string, string, string, error) {
	return scriptTransformer.CompileJavascriptFile(scriptPath, gEngineOptions.ScriptAutomaticHeader, gEnableDebugger)
}

func executeScript(ctx progpAPI.JsContext, scriptPath string, onCompiledSuccess func()) *progpAPI.JsErrorMessage {
	// Transform typescript file (and others supported types) as plain javascript.
	// It big a big file with all the requirements.
	//
	scriptContent, scriptOrigin, sourceMap, err := compileScript(scriptPath)

	// If ko, the error message has already been displayed.
	// Then we only have to exit.
	//
	if err != nil {
		if gEngineOptions.OnScriptCompilationError != nil {
			if gEngineOptions.OnScriptCompilationError(scriptPath, err) {
				return nil
			}
		}

		os.Exit(1)
	}

	onCompiledSuccess()
	return ctx.ExecuteScript(scriptContent, scriptOrigin, scriptPath, sourceMap)
}

// bootstrapWithOptions initialize the engine and execute a startup script.
// If the script path is blank, then no script is executed.
// In all case the engine is initialized.
func bootstrapWithOptions(options *EngineOptions, installMods func()) {
	if gIsBootstrapped {
		return
	}

	gIsBootstrapped = true
	gEngineOptions = options
	gEnableDebugger = options.MustDebug

	// Configure things for the core functionalities.
	configureCore()

	// Get instance of the engine or panic if not found.
	//
	// This instance is registered by "scriptEngine.progpV8" if linked to the source.
	// If not will load progpV8 as a plugin from the file which path is "../plugins/progpV8.so".
	//
	scriptEngine := resolveScriptEngine(options.ScriptEngineName, options.PluginsDir)
	gDefaultScriptEngine = scriptEngine

	// Warning the engine must be defined before calling.
	installMods()

	// Get the function registry and declare all the function to the script engine implementation.
	// Will create dynamic function, or update the compiled code if env variable PROGPV8_DIR
	// points to the source dir of "scriptEngine.progpV8".
	//
	GenerateSourceCode(options.ProgpV8EngineProjectDir)

	if options.OnRuntimeError != nil {
		scriptEngine.SetRuntimeErrorHandler(options.OnRuntimeError)
	}

	scriptEngine.SetScriptTerminatedHandler(options.OnScriptTerminated)
	scriptEngine.SetAllowedFunctionsChecker(options.OnCheckingAllowedFunction)

	// Allows the engine to initialize himself.
	scriptEngine.Start()

	progpAPI.SetScriptFileExecutor(executeScript)
	progpAPI.SetScriptFileCompiler(compileScript)

	gSignalHandler = onProgpJsSignal

	// Allows closing resources correctly and
	// avoid some errors which can occurs before exiting.
	//
	runtime.GC()
}

var gSignalHandler progpAPI.ListenProgpSignalF

func EmitProgpSignal(ctx progpAPI.JsContext, signal string, data string) error {
	if gSignalHandler != nil {
		return gSignalHandler(ctx, signal, data)
	}

	return nil
}

func WaitEnd(forceEnd bool) {
	if !forceEnd {
		// Wait until all background task finished.
		// A background task can be a webserver is list for call.
		// In this case the tasks never ends until the server is stopped.
		//
		progpAPI.WaitTasksEnd()
	}

	progpAPI.ForceExitingVM()
}

// ExecuteScriptFile is like ExecuteScript but allows using a file (which can be typescript).
func ExecuteScriptFile(scriptPath string, securityGroup string, mustDebug bool, onCompiledSuccess func()) *progpAPI.JsErrorMessage {
	ctx := gDefaultScriptEngine.CreateNewScriptContext(securityGroup, mustDebug)

	if mustDebug {
		gDefaultScriptEngine.WaitDebuggerReady()
	}

	return ctx.ExecuteScriptFile(scriptPath, onCompiledSuccess)
}

func GenerateSourceCode(progpV8EngineProjectDir string) {
	// If the directory is provided then build an optimized version of the sources
	// which avoid using reflection. It's much more fast!
	//
	// Without that, it uses "draft functions" which use reflection and are very slow.

	mustUseDynamicMode := progpV8EngineProjectDir == ""

	if !mustUseDynamicMode {
		if !path.IsAbs(progpV8EngineProjectDir) {
			cwd, _ := os.Getwd()
			progpV8EngineProjectDir = path.Join(cwd, progpV8EngineProjectDir)
		}
	}

	fctRegistry := progpAPI.GetFunctionRegistry()
	fctRegistry.EnableDynamicMode(mustUseDynamicMode)

	if !mustUseDynamicMode {
		codeGen := codegen.NewProgpV8Codegen()
		codeGen.GenerateCode(progpV8EngineProjectDir)
	}
}

func configureCore() {
	// Will allows to translate error message from plain javascript to typescript.
	// This by using a sourcemap to decode.
	//
	progpAPI.SetErrorTranslator(func(message *progpAPI.JsErrorMessage) {
		executingScript := message.ScriptPath

		sourceMap := []byte(message.SourceMap)

		if len(sourceMap) == 0 {
			txt, err := os.ReadFile(executingScript + ".map")
			if err != nil {
				return
			}
			sourceMap = txt
		}

		reader, err := sourcemap.Parse("", sourceMap)
		if err != nil {
			return
		}

		for offset, frame := range message.StackTraceFrames {
			file, _, line, col, ok := reader.Source(frame.Line, frame.Column)

			if ok {
				frame.Column = col
				frame.Line = line
				frame.Source = file

				message.StackTraceFrames[offset] = frame
			}
		}
	})

	// Allows to embed modules into go code.
	//
	// A call will be done for each module which aren't found.
	// Here the result must be one plain javascript uniq file.
	// Which mean that the stored module must be compiled.
	//
	scriptTransformer.SetJavascriptModuleResolver(resolveMissingJavascriptModule)
}
