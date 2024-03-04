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
	"embed"
	"github.com/progpjs/progpAPI/v2"
	"github.com/progpjs/progpjs/v2/scriptTransformer"
	"log"
)

func GetFunctionRegistry() *progpAPI.FunctionRegistry {
	return progpAPI.GetFunctionRegistry()
}

func ReadEmbeddedFile(fs embed.FS, innerPath string) string {
	b, err := fs.ReadFile(innerPath)
	if err != nil {
		log.Print("Error loading embedded resource ", innerPath)
		return ""
	}
	return string(b)
}

func ReturnEmbeddedTypescriptModule(fs embed.FS, innerPath string) JavascriptModuleProviderF {
	return func(modName string) (content string, loader JsResourceLoader) {
		loader = JsLoaderTS
		content = ReadEmbeddedFile(fs, innerPath)
		return
	}
}

var gSignalListeners []progpAPI.ListenProgpSignalF

func AddSignalListener(listener progpAPI.ListenProgpSignalF) {
	gSignalListeners = append(gSignalListeners, listener)
}

func onProgpJsSignal(ctx progpAPI.JsContext, signal string, data string) error {
	for _, listener := range gSignalListeners {
		err := listener(ctx, signal, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCacheDir(searchFrom string, createDir bool) string {
	return scriptTransformer.GetCacheDir(searchFrom, createDir)
}
