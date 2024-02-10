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

package libProgpScripts

import (
	"github.com/go-sourcemap/sourcemap"
	"github.com/progpjs/libProgpScripts/scriptTransformer"
	"github.com/progpjs/progpAPI"
	"github.com/progpjs/progpAPI/codegen"
	"os"
	"path"
	"plugin"
	"runtime"
	"strings"
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

func loadGoPlugins() {
	if gGoPluginAreLoaded {
		return
	}
	gGoPluginAreLoaded = true

	cwd, _ := os.Getwd()
	pluginDir := path.Join(cwd, "..", "plugins")
	loadGoPlugin(path.Join(pluginDir, "progpV8.so"))
}

//endregion

func getScriptEngine() progpAPI.ScriptEngine {
	// Currently there is only one engine.
	const engineName = "progpV8"

	if gScriptEngine != nil {
		return gScriptEngine
	}

	gScriptEngine = progpAPI.GetScriptEngine(engineName)
	if gScriptEngine != nil {
		return gScriptEngine
	}

	loadGoPlugins()
	gScriptEngine = progpAPI.GetScriptEngine(engineName)
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

var gJavascriptModuleProviders = make(map[string]JavascriptModuleProviderF)

type JavascriptModuleProviderF func(resourcePath string) (content string, loader JsResourceLoader)

//endregion

func StartupEngine(scriptPath string, launchDebugger bool) {
	// Get the function registry and declare all the function to the script engine implementation.
	// Will create dynamic function, or update the compiled code if env variable PROGPV8_DIR
	// points to the source dir of "scriptEngine.progpV8".
	//
	exportExposedFunctions()

	// Get instance of the engine or panic if not found.
	//
	// This instance is registered by "scriptEngine.progpV8" if linked to the source.
	// If not will load progpV8 as a plugin from the file which path is "../plugins/progpV8.so".
	//
	scriptEngine := getScriptEngine()

	// Configure things for to the engine.
	configureScriptEngine()

	// Transform typescript file (and others supported types) as plain javascript.
	// It big a big file with all the requirements.
	//
	scriptContent, scriptOrigin, isOk := scriptTransformer.CompileJavascriptFile(scriptPath)

	// If ko, the error message has already been displayed.
	// Then we only have to exit.
	//
	if !isOk {
		os.Exit(1)
	}

	scriptEngine.Start()

	if launchDebugger {
		scriptEngine.WaitDebuggerReady()
	}

	progpAPI.ExecuteScriptContent(scriptContent, scriptOrigin, scriptEngine)

	// Allows closing resources correctly and
	// avoid some errors which can occurs before exiting.
	//
	runtime.GC()
}

func exportExposedFunctions() {
	// If the directory is provided then build an optimized version of the sources
	// which avoid using reflection. It's much more fast!
	//
	// Without that, it uses "draft functions" which use reflection and are very slow.
	//
	progpV8Dir := strings.Trim(os.Getenv("PROGPV8_DIR"), " ")

	mustUseDynamicMode := progpV8Dir == ""
	fctRegistry := progpAPI.GetFunctionRegistry()
	fctRegistry.EnableDynamicMode(mustUseDynamicMode)

	if !mustUseDynamicMode {
		codeGen := codegen.NewProgpV8Codegen()
		codeGen.GenerateCode(progpV8Dir)
	}
}

func configureScriptEngine() {
	// Will allows to translate error message from plain javascript to typescript.
	// This by using a sourcemap to decode.
	//
	progpAPI.SetErrorTranslator(func(message *progpAPI.ScriptErrorMessage) {
		executingScript := message.ScriptPath

		sourceMapFileContent, err := os.ReadFile(executingScript + ".map")
		if err != nil {
			return
		}

		reader, err := sourcemap.Parse("", sourceMapFileContent)
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
