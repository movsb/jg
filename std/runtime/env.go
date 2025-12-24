package jg_runtime

import (
	"os"
	"runtime"
)

var Map = map[string]any{
	`os`:   runtime.GOOS,
	`arch`: runtime.GOARCH,
	`args`: os.Args,
}
