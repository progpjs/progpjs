package scriptTransformer

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path"
)

func calcMd5(text string) string {
	h := md5.New()
	//_, _ = io.WriteString(h, text)
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func FileExists(path string) bool {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if !stat.IsDir() {
		return true
	}

	return false
}

func DirExists(path string) bool {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if stat.IsDir() {
		return true
	}

	return false
}

func SearchNodeModulesDir(searchFrom string) string {
	dirPath := path.Join(searchFrom, "node_modules")

	if DirExists(dirPath) {
		return dirPath
	}

	parent := path.Dir(searchFrom)
	if (parent == "") || (parent == searchFrom) {
		return ""
	}

	return SearchNodeModulesDir(parent)
}

func GetCacheDir(searchFrom string, createDir bool) string {
	if gCacheDir == "" {
		nmd := SearchNodeModulesDir(searchFrom)

		const cacheDirName = ".progpjs"

		if nmd != "" {
			gCacheDir = path.Join(path.Dir(nmd), cacheDirName)
		} else {
			gCacheDir = path.Join(searchFrom, cacheDirName)
		}
	}

	if createDir {
		if !gIsCacheDirCreate {
			gIsCacheDirCreate = true
			_ = os.MkdirAll(gCacheDir, os.ModePerm)
		}

	}

	return gCacheDir
}

func GetCompileCacheDir(searchFrom string, createDir bool) string {
	return path.Join(GetCacheDir(searchFrom, createDir), "build")
}

// SearchModuleInNodeModules returns the module path inside node_modules directory.
// It recurse in the upper directory if not found in the current one.
func SearchModuleInNodeModules(modName string, callerDir string) string {
	var foundPath string

	// Here add sub_dir in order to match the case where the caller dir is node_modules.
	//
	WalkNodeModules(path.Join(callerDir, "/sub_dir"), func(nodeModulesDir string) bool {
		basePath := path.Join(nodeModulesDir, modName)
		found := searchFileFromBase(basePath)

		if found != "" {
			foundPath = found
			return true
		}

		return false
	})

	return foundPath
}

// WalkNodeModules allows a function to be called on each node_modules dir until
// this function return true, going up recursively on the node_modules  in the parent dir.
func WalkNodeModules(startDir string, handler func(dirPath string) bool) {
	dirToTest := path.Join(startDir, "node_modules")

	info, err := os.Stat(dirToTest)

	if err == nil && info.IsDir() {
		if handler(dirToTest) {
			return
		}
	}

	parentDir := path.Dir(startDir)
	if parentDir == startDir {
		return
	}

	WalkNodeModules(parentDir, handler)
}

var gCacheDir string
var gIsCacheDirCreate bool
