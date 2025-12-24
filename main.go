package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/dop251/goja"
	loop "github.com/movsb/jg/runtime"
	jg_archive "github.com/movsb/jg/std/archive"
	jg_fs "github.com/movsb/jg/std/fs"
	jg_http "github.com/movsb/jg/std/http"
	jg_exec "github.com/movsb/jg/std/os/exec"
	jg_path "github.com/movsb/jg/std/path"
	jg_runtime "github.com/movsb/jg/std/runtime"
)

func main() {
	rt, err := loop.NewRuntime(context.Background())
	if err != nil {
		panic(err)
	}
	rt.Run(func(rt *goja.Runtime) {
		rt.SetFieldNameMapper(goja.UncapFieldNameMapper())
	})
	rt.Run(func(rt *goja.Runtime) {
		rt.Set(`runtime`, jg_runtime.Map)
		rt.Set(`http`, jg_http.Methods)
		rt.Set(`fs`, jg_fs.Methods)
		rt.Set(`path`, &jg_path.Path{})
		rt.Set(`exec`, jg_exec.Methods)
		rt.Set(`$`, jg_exec.Methods[`$`])
		rt.Set(`archive`, &jg_archive.Methods)
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

	output, err := rt.Execute(context.Background(), string(all))
	if err != nil {
		log.Fatalln(err)
	}

	if promise, ok := output.(*goja.Promise); ok {
		rt.WaitForPromise(promise)
		if promise.State() == goja.PromiseStateRejected {
			log.Println(promise.Result())
		}
	}
}
