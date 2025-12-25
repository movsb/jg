package loop

import (
	"fmt"

	"github.com/dop251/goja"
)

type Runtime interface {
	RunAsync(fn func(rt *goja.Runtime))
}

func GetRunAsync(vm *goja.Runtime) func(func(vm *goja.Runtime)) {
	return vm.Get(`__loop`).Export().(Runtime).RunAsync
}

func CreatePromise(vm *goja.Runtime) (_ *goja.Promise, _, _ func(value any)) {
	promise, resolve, reject := vm.NewPromise()
	return promise,
		func(value any) {
			GetRunAsync(vm)(func(vm *goja.Runtime) {
				resolve(value)
			})
		},
		func(value any) {
			GetRunAsync(vm)(func(vm *goja.Runtime) {
				reject(value)
			})
		}
}

func Panic(vm *goja.Runtime, format string, args ...any) {
	panic(vm.ToValue(fmt.Errorf(format, args...)))
}
