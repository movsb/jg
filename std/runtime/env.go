package jg_runtime

import (
	"runtime"
)

var Map = map[string]any{
	`os`:   runtime.GOOS,
	`arch`: runtime.GOARCH,
}
