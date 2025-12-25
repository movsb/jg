package jg_runtime

import (
	"fmt"
	"os"
	"runtime"
)

var Map = map[string]any{
	`os`:    runtime.GOOS,
	`arch`:  runtime.GOARCH,
	`args`:  os.Args,
	`panic`: _panic,
}

func _panic(v ...any) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}
