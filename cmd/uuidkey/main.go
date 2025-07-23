package main

import (
	"fmt"
	"os"
	"runtime"
)

// Build variables set by ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// Runtime variables
var (
	goVersion = runtime.Version()
	platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// osExit is a variable to allow mocking in tests
var osExit = os.Exit

func main() {
	if err := Execute(); err != nil {
		osExit(1)
	}
}