package jg_http

import (
	"io"
	"net/http"
	"runtime"

	"github.com/dop251/goja"
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

func (r *Response) Text() (string, error) {
	b, err := io.ReadAll(r.r.Body)
	return string(b), err
}

func (r *Response) Blob(args goja.FunctionCall, vm *goja.Runtime) goja.Value {
	b, err := io.ReadAll(r.r.Body)
	if err != nil {
		panic(vm.ToValue(err))
	}
	return vm.ToValue(vm.NewArrayBuffer(b))
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
	rsp, err := http.Get(u)
	if err != nil {
		panic(vm.ToValue(err))
	}
	rsp2 := NewResponse(rsp)
	runtime.AddCleanup(rsp2, func(int) {
		rsp2.r.Body.Close()
	}, 0)
	return vm.ToValue(rsp2)
}
