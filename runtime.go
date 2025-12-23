package main

import (
	"context"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/eventloop"
)

// JavaScript 运行时，必要时移出去作为公共模块。
type Runtime struct {
	ev *eventloop.EventLoop
}

func NewRuntime(ctx context.Context, libs ...[]byte) (*Runtime, error) {
	ev := eventloop.NewEventLoop()
	ev.Start()

	run(ev, func(r *goja.Runtime) {
		console.Enable(r)
		for _, lib := range libs {
			_, err := r.RunString(string(lib))
			if err != nil {
				panic(err)
			}
		}
	})

	return &Runtime{ev: ev}, nil

}

func run(ev *eventloop.EventLoop, fn func(r *goja.Runtime)) {
	wait := make(chan struct{})
	ev.RunOnLoop(func(r *goja.Runtime) {
		defer close(wait)
		fn(r)
	})
	<-wait
}

type Argument struct {
	Name  string
	Value any
}

func (r *Runtime) Run(fn func(rt *goja.Runtime)) {
	run(r.ev, fn)
}

// arguments 被设置到全局，执行完成后删除。
func (r *Runtime) Execute(ctx context.Context, script string, arguments ...Argument) (_ goja.Value, outErr error) {
	var val goja.Value
	var err error

	run(r.ev, func(r *goja.Runtime) {
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
