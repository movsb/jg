package jg_exec

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dop251/goja"
	loop "github.com/movsb/jg/runtime"
	"github.com/movsb/jg/utils"
)

var Methods = map[string]any{
	`Command`: command,
}

type Command struct {
	underlying *exec.Cmd
}

func (c *Command) Run(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	async := loop.GetRunAsync(vm)
	promise, resolve, reject := vm.NewPromise()

	go func() {
		err := c.underlying.Run()
		if err != nil {
			async(func(vm *goja.Runtime) {
				reject(vm.ToValue(fmt.Errorf(`运行命令时出错：%w`, err)))
			})
			return
		}
		async(func(vm *goja.Runtime) {
			resolve(goja.Undefined())
		})
	}()

	return vm.ToValue(promise)
}

func (c *Command) UseStd(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if in := utils.MustBeBoolean(call.Argument(0), vm); in {
		c.underlying.Stdin = os.Stdin
	}
	if out := utils.MustBeBoolean(call.Argument(1), vm); out {
		c.underlying.Stdout = os.Stdout
	}
	if err := utils.MustBeBoolean(call.Argument(2), vm); err {
		c.underlying.Stderr = os.Stderr
	}
	return goja.Undefined()
}

func command(call goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	cmd := exec.Command(call.Argument(0).String())
	for i := 1; i < len(call.Arguments); i++ {
		utils.MustBeString(call.Argument(i), vm)
		cmd.Args = append(cmd.Args, call.Argument(i).String())
	}

	myCmd := &Command{
		underlying: cmd,
	}
	obj := vm.ToValue(myCmd).(*goja.Object)
	obj.SetPrototype(call.This.Prototype())

	return obj
}
