package jg_fs

import (
	"fmt"
	"io"
	"os"

	"github.com/dop251/goja"
	loop "github.com/movsb/jg/runtime"
	"github.com/movsb/jg/utils"
)

var Methods = map[string]any{
	`saveToFile`: saveToFile,
}

func saveToFile(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	filePath := utils.MustBeString(call.Argument(0), vm)
	from, ok := call.Argument(1).Export().(io.Reader)
	if !ok {
		panic(vm.ToValue(`第二个参数未实现 Reader 接口`))
	}

	async := loop.GetRunAsync(vm)
	promise, resolve, reject := vm.NewPromise()

	go func() {
		fp, err := os.Create(filePath)
		if err != nil {
			async(func(vm *goja.Runtime) {
				reject(vm.ToValue(fmt.Errorf(`创建文件失败：%w`, err)))
			})
			return
		}

		n, err := io.Copy(fp, from)
		if err != nil {
			fp.Close()
			os.Remove(filePath)
			async(func(vm *goja.Runtime) {
				reject(vm.ToValue(fmt.Errorf(`写文件失败：%w`, err)))
			})
		}

		if err := fp.Close(); err != nil {
			os.Remove(filePath)
			async(func(vm *goja.Runtime) {
				reject(vm.ToValue(fmt.Errorf(`写文件失败：%w`, err)))
			})
		}

		resolve(n)
	}()

	return vm.ToValue(promise)
}
