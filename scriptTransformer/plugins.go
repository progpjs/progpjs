package scriptTransformer

import (
	"errors"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"path"
)

var gPlugins []api.Plugin

func getPlugins() []api.Plugin {
	if gPlugins == nil {
		gPlugins = []api.Plugin{pluginResolveMissingDependency()}
	}

	return gPlugins
}

func searchFileFromBase(filePath string) string {
	allExt := []string{".tsx", ".ts", ".js"}

	for _, ext := range allExt {
		f := filePath + ext
		state, err := os.Stat(f)

		if (err == nil) && !state.IsDir() {
			return f
		}
	}

	filePath = path.Join(filePath, "index")

	for _, ext := range allExt {
		f := filePath + ext
		state, err := os.Stat(f)

		if (err == nil) && !state.IsDir() {
			return f
		}
	}

	return ""
}

// pluginResolveMissingDependency allows requesting a javascript
// module which is located elsewhere or which is provider by an external source.
func pluginResolveMissingDependency() api.Plugin {
	return api.Plugin{
		Name: "progpResolveMissingDependency",

		Setup: func(build api.PluginBuild) {
			//region Resolve - search the file location or add him a namespace

			// To know: here the resolver don't provide the file content but
			// only his path. It's why we must use a resolver before getting
			// to the loader step.

			addToProgpNS := func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				// Search in the node_modules hierarchy.
				// This step is required, since esbuild only look for ".js" files.
				// Here with this patch we can use ".ts", ".tsx" and ".jsx" files.
				//
				foundPath := SearchModuleInNodeModules(args.Path, args.ResolveDir)
				//
				if foundPath != "" {
					return api.OnResolveResult{
						Path:      foundPath,
						Namespace: "file",
					}, nil
				}

				return api.OnResolveResult{
					Path:       args.Path,
					Namespace:  "progp",
					PluginData: args.ResolveDir,
				}, nil
			}

			filters := []string{
				// progp packages
				`^@progp/`,
				`^progp:`,

				// Allows flagging explicitly embedded files
				`^embedded:`,

				// Will allows replacing by our own version
				`^react$`,

				// Node.js namespace
				`^node:`,

				// Node.js special packages
				`^assert$`, `^path$`, `^fs$`, `^os$`, `^process$`, `^stream$`, `^test$`,
			}

			for _, filter := range filters {
				build.OnResolve(api.OnResolveOptions{Filter: filter}, addToProgpNS)
			}

			//endregion

			//region Loader - load file from a namespace

			// Will call the provider with the searched module.
			onLoad := func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				modName := args.Path

				if gJavascriptModuleResolver != nil {
					asText, loader, ok := gJavascriptModuleResolver(modName)

					if ok {
						return api.OnLoadResult{
							Contents: &asText,
							Loader:   api.Loader(loader),
						}, nil
					}
				}

				return api.OnLoadResult{}, errors.New("can't found dependency")
			}

			// Takes content from Go embedded file.
			// Allows catching things like 		import "progp:core".
			// and								import "embedded:myscript".
			//
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "progp"}, onLoad)
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "embedded"}, onLoad)

			// For Node.js compatibility.
			// Allows catching things like 		import "node:path".
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "node"}, onLoad)

			//endregion
		},
	}
}

type JavascriptModuleResolverF func(resourcePath string) (content string, loaderToUse uint16, isFound bool)

var gJavascriptModuleResolver JavascriptModuleResolverF
