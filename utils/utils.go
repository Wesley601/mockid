package utils

import (
	"os"
	"path"
	"runtime"
)

func Must[T any](v T, err error) T {
	if err == nil {
		return v
	}
	panic(err)
}

func Env(key, d string) string {
	value := os.Getenv(key)
	if value == "" {
		return d
	}

	return value
}

// for some reason the go test runner run the test from the test dir instead of the root dir
// this may cause a series of inconsistences.
// so this function need to be call on all init test method nested from the root location
// ex:
//
//	func init() {
//		utils.SetToRoot()
//	}
func SetToRoot(level string) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), level)
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
