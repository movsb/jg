package jg_fs

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/dop251/goja"
	"github.com/movsb/jg/runtime/loop"
	"github.com/movsb/jg/utils"
)

var Methods = map[string]any{
	`saveToFile`: saveToFile,

	`fileExists`: fileExists,

	`sha256`: sha256sum,
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

func fileExists(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	filePath := utils.MustBeString(call.Argument(0), vm)

	var types string

	if len(call.Arguments) >= 2 {
		types = utils.MustBeString(call.Argument(1), vm)
	} else {
		types = `fdls`
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return vm.ToValue(false)
		}
		panic(vm.ToValue(fmt.Errorf(`判断存在时出错：%s: %w`, filePath, err)))
	}

	orMatch := false
	andMatch := true

	for _, t := range types {
		switch t {
		case 'f':
			orMatch = orMatch || stat.Mode().IsRegular()
		case 'd':
			orMatch = orMatch || stat.IsDir()
		case 'l':
			orMatch = orMatch || (stat.Mode().Type()&fs.ModeSymlink > 0)
		case 's':
			orMatch = orMatch || (stat.Mode().Type()&fs.ModeSocket > 0)

		case 'x':
			andMatch = andMatch && (stat.Mode().Perm()&0b001 > 0)
		case 'r':
			andMatch = andMatch && (stat.Mode().Perm()&0b010 > 0)
		case 'w':
			andMatch = andMatch && (stat.Mode().Perm()&0b100 > 0)
		}
	}

	return vm.ToValue(orMatch && andMatch)
}
