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
	"os"
	"path"
)

func CreateDefaultEngineOptions(enableDebug bool) *EngineOptions {
	options := &EngineOptions{}
	cwd, _ := os.Getwd()

	options.MustDebug = enableDebug
	options.PluginsDir = path.Join(cwd, "..", "..", "_plugins")
	options.ProgpV8EngineProjectDir = os.Getenv("PROGPV8_DIR")

	// Optional, allows selecting the engine when more than one is available.
	options.ScriptEngineName = "progpV8"

	return options
}

func Bootstrap(scriptPath string, enableDebug bool, options *EngineOptions) {
	if options == nil {
		options = CreateDefaultEngineOptions(enableDebug)
	}

	if scriptPath != "" && !path.IsAbs(scriptPath) {
		cwd, _ := os.Getwd()
		scriptPath = path.Join(cwd, scriptPath)
	}

	// bootstrapWithOptions the engine.
	bootstrapWithOptions(options)

	// Execute our script.
	//
	// The current thread block until the script has totally terminated to execute.
	// If it's not what you want, then add "go " before in order to create a new thread.
	//
	// Here we set the security group to "admin". The meaning is related to options.OnCheckingAllowedFunction
	// and the rules you put here.
	//
	scriptErr := ExecuteScriptFile(scriptPath, "admin", options.MustDebug)

	// Will wait until all background tasks termine and dispose the script engine.
	// A background task is for exemple a web server listening a port.
	// In this case, it's never ends.
	//
	// Calling this function is important, since without that the app exit immediately.
	//
	WaitEnd(scriptErr != nil)
}
