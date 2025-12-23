package jg_exec

import (
	"os/exec"

	"github.com/dop251/goja"
	"github.com/movsb/jg/utils"
)

type Exec struct{}

func (*Exec) Command(call goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	cmd := exec.Command(call.Argument(0).String())
	for i := 1; i < len(call.Arguments); i++ {
		utils.MustBeString(call.Argument(i), vm)
		cmd.Args = append(cmd.Args, call.Argument(i).String())
	}
	obj := vm.ToValue(cmd).(*goja.Object)
	obj.SetPrototype(call.This.Prototype())
	return obj
}
