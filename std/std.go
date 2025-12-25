package std

import (
	"github.com/dop251/goja"
	jg_archive "github.com/movsb/jg/std/archive"
	jg_fs "github.com/movsb/jg/std/fs"
	jg_http "github.com/movsb/jg/std/http"
	jg_exec "github.com/movsb/jg/std/os/exec"
	jg_path "github.com/movsb/jg/std/path"
)

func Init(rt *goja.Runtime) {
	rt.Set(`http`, jg_http.Methods)
	rt.Set(`fs`, jg_fs.Methods)
	rt.Set(`path`, &jg_path.Path{})
	rt.Set(`exec`, jg_exec.Methods)
	rt.Set(`$`, jg_exec.Methods[`$`])
	rt.Set(`archive`, &jg_archive.Methods)
}
