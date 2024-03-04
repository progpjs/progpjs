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
	"github.com/progpjs/progpAPI/v2"
	"os"
	"path"
	"sync"
)

func DefaultBootstrapOptions() *EngineOptions {
	options := &EngineOptions{}
	cwd, _ := os.Getwd()

	options.PluginsDir = path.Join(cwd, "..", "..", "_plugins")
	options.ProgpV8EngineProjectDir = os.Getenv("PROGPV8_DIR")

	if (options.ProgpV8EngineProjectDir != "") && !path.IsAbs(options.ProgpV8EngineProjectDir) {
		options.ProgpV8EngineProjectDir = path.Join(cwd, options.ProgpV8EngineProjectDir)
	}

	// Optional, allows selecting the engine when more than one is available.
	options.ScriptEngineName = "progpV8"

	options.OnScriptCompilationError = func(scriptPath string, err error) bool {
		print(err.Error())
		return false
	}

	return options
}

func Bootstrap(scriptPath string, enableDebug bool, options *EngineOptions, installMods func()) BootstrapExitAwaiterF {
	if options == nil {
		options = DefaultBootstrapOptions()
	}

	options.MustDebug = enableDebug

	if scriptPath != "" && !path.IsAbs(scriptPath) {
		cwd, _ := os.Getwd()
		scriptPath = path.Join(cwd, scriptPath)
	}

	// bootstrapWithOptions the engine.
	bootstrapWithOptions(options, installMods)

	// Execute our script.
	// Here we execute the script in a parallel thread in order to don't block the caller.
	// This allows to make things more simple when embedding the engine.
	//
	var scriptErr *progpAPI.JsErrorMessage
	var mutex sync.Mutex
	mutex.Lock()
	//
	go func() {
		isUnlocked := false

		// Here we set the security group to "admin". The meaning is related to options.OnCheckingAllowedFunction
		// and the rules you put here.
		//
		scriptErr = ExecuteScriptFile(scriptPath, "admin", options.MustDebug, func() {
			// Unlock if script compilation is ok.
			isUnlocked = true
			mutex.Unlock()
		})

		if !isUnlocked {
			// Unlock if error.
			mutex.Lock()
		}
	}()

	mutex.Lock()

	return func() {
		// Will wait until all background tasks terminate and dispose the script engine.
		// A background task is for example a web server listening a port.
		// In this case, it's never ends.
		//
		// Calling this function is important, since without that the app exit immediately.
		//
		WaitEnd(scriptErr != nil)
	}
}

type BootstrapExitAwaiterF func()
