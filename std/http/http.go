package jg_http

import (
	"io"
	"net/http"
	"runtime"

	"github.com/dop251/goja"
	loop "github.com/movsb/jg/runtime"
	"github.com/movsb/jg/utils"
)

var Methods = map[string]any{
	`get`: get,
}

type Response struct {
	r *http.Response

	StatusCode int
	Status     string
}

func (r *Response) Reader(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	return vm.ToValue(r.r.Body)
}

func (r *Response) Text(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	async := loop.GetRunAsync(vm)
	promise, resolve, reject := vm.NewPromise()

	go func() {
		b, err := io.ReadAll(r.r.Body)
		if err != nil {
			async(func(vm *goja.Runtime) {
				reject(vm.ToValue(err))
			})
			return
		}
		async(func(vm *goja.Runtime) {
			resolve(string(b))
		})
	}()

	return vm.ToValue(promise)
}

func (r *Response) Blob(args goja.FunctionCall, vm *goja.Runtime) goja.Value {
	async := loop.GetRunAsync(vm)
	promise, resolve, reject := vm.NewPromise()

	go func() {
		b, err := io.ReadAll(r.r.Body)
		if err != nil {
			async(func(vm *goja.Runtime) {
				reject(vm.ToValue(err))
			})
			return
		}
		async(func(vm *goja.Runtime) {
			resolve(vm.NewArrayBuffer(b))
		})
	}()

	return vm.ToValue(promise)
}

func NewResponse(r *http.Response) *Response {
	rr := &Response{
		r:          r,
		Status:     r.Status,
		StatusCode: r.StatusCode,
	}
	return rr
}

func get(args goja.FunctionCall, vm *goja.Runtime) goja.Value {
	u := utils.MustBeString(args.Argument(0), vm)

	loop := loop.GetRunAsync(vm)
	promise, resolve, reject := vm.NewPromise()

	go func() {
		rsp, err := http.Get(u)
		if err != nil {
			loop(func(vm *goja.Runtime) {
				reject(vm.ToValue(err))
			})
			return
		}
		rsp2 := NewResponse(rsp)
		runtime.AddCleanup(rsp2, func(int) {
			rsp2.r.Body.Close()
		}, 0)
		loop(func(vm *goja.Runtime) {
			resolve(vm.ToValue(rsp2))
		})
	}()

	return vm.ToValue(promise)
}
