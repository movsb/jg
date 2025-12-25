package jr

import (
	"context"
	"fmt"
	"os"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/movsb/jg/runtime/loop"
	"github.com/movsb/jg/utils"
)

// JavaScript 运行时，必要时移出去作为公共模块。
type Runtime struct {
	ev *eventloop.EventLoop
}

var _ loop.Runtime = (*Runtime)(nil)

func MustNewRuntime(ctx context.Context, options ...Option) *Runtime {
	return utils.Must1(NewRuntime(ctx, options...))
}

func NewRuntime(ctx context.Context, options ...Option) (*Runtime, error) {
	rt := &Runtime{ev: eventloop.NewEventLoop()}

	rt.Run(func(vm *goja.Runtime) {
		vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
		console.Enable(vm)
		vm.Set(`__loop`, rt)
		vm.Set(`runtime`, _common)
		vm.Set(`panic`, _common[`panic`])
	})

	for _, opt := range options {
		opt(rt)
	}

	return rt, nil
}

type Argument struct {
	Name  string
	Value any
}

func (r *Runtime) Run(fn func(rt *goja.Runtime)) {
	r.ev.Run(fn)
}

func (r *Runtime) RunAsync(fn func(rt *goja.Runtime)) {
	r.ev.RunOnLoop(fn)
}

// arguments 被设置到全局，执行完成后删除。
func (r *Runtime) Execute(ctx context.Context, script string, arguments ...Argument) (_ any, outErr error) {
	var val goja.Value
	var err error

	r.Run(func(r *goja.Runtime) {
		for _, arg := range arguments {
			if err := r.Set(arg.Name, arg.Value); err != nil {
				panic(err)
			}
		}

		defer func() {
			for _, arg := range arguments {
				r.GlobalObject().Delete(arg.Name)
			}
		}()
		val, err = r.RunString(script)
	})

	if err != nil {
		return nil, err
	}

	return val.Export(), nil
}

func (r *Runtime) ExecuteFile(ctx context.Context, file string) (any, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	output, err := r.Execute(ctx, string(b))
	if err != nil {
		return nil, err
	}

	if promise, ok := output.(*goja.Promise); ok {
		r.waitForPromise(promise)
		if promise.State() == goja.PromiseStateRejected {
			return nil, fmt.Errorf(`%v`, promise.Result().Export())
		}
		return promise.Result().Export(), nil
	}

	return output, nil
}

func (r *Runtime) waitForPromise(p *goja.Promise) {
	var timer *eventloop.Interval
	timer = r.ev.SetInterval(func(rt *goja.Runtime) {
		if p.State() != goja.PromiseStatePending {
			r.ev.ClearInterval(timer)
		}
	}, 100)
	r.ev.Run(func(r *goja.Runtime) {})
}
