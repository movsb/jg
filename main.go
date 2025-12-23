package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/dop251/goja"
	jg_http "github.com/movsb/jg/std/http"
	jg_exec "github.com/movsb/jg/std/os/exec"
	jg_path "github.com/movsb/jg/std/path"
	jg_runtime "github.com/movsb/jg/std/runtime"
)

func main() {
	rt, err := NewRuntime(context.Background())
	if err != nil {
		panic(err)
	}
	rt.Run(func(rt *goja.Runtime) {
		rt.SetFieldNameMapper(goja.UncapFieldNameMapper())
	})
	rt.Run(func(rt *goja.Runtime) {
		rt.Set(`runtime`, &jg_runtime.Map)
		rt.Set(`http`, &jg_http.Methods)
		rt.Set(`path`, &jg_path.Path{})
		rt.Set(`exec`, &jg_exec.Exec{})
	})

	fp, err := os.Open(`main.js`)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	all, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}

	val, err := rt.Execute(context.Background(), string(all))
	if err != nil {
		log.Fatalln(err)
	}
	// log.Println(val.Export())
	_ = val
}
