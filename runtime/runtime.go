package loop

import (
	"context"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/eventloop"
)

func GetRunAsync(vm *goja.Runtime) func(func(vm *goja.Runtime)) {
	rt := vm.Get(`__loop`).Export().(*Runtime)
	return rt.RunAsync
}

// JavaScript 运行时，必要时移出去作为公共模块。
type Runtime struct {
	ev *eventloop.EventLoop
}

func NewRuntime(ctx context.Context, libs ...[]byte) (*Runtime, error) {
	rt := &Runtime{ev: eventloop.NewEventLoop()}

	rt.Run(func(r *goja.Runtime) {
		console.Enable(r)
		for _, lib := range libs {
			_, err := r.RunString(string(lib))
			if err != nil {
				panic(err)
			}
		}
	})

	rt.Run(func(r *goja.Runtime) {
		r.Set(`__loop`, rt)
	})

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
func (r *Runtime) Execute(ctx context.Context, script string, arguments ...Argument) (_ goja.Value, outErr error) {
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

	return val, err
}

func (r *Runtime) WaitForPromise(p *goja.Promise) {
	var timer *eventloop.Interval
	timer = r.ev.SetInterval(func(rt *goja.Runtime) {
		if p.State() != goja.PromiseStatePending {
			r.ev.ClearInterval(timer)
		}
	}, 100)
	r.ev.Run(func(r *goja.Runtime) {})
}
