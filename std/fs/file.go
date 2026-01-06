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
	`mkDir`:    mkDir,
	`mkDirAll`: mkDirAll,

	`saveToFile`: saveToFile,

	`exists`:     exists,
	`fileExists`: fileExists,
	`dirExists`:  dirExists,

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

func _exists(vm *goja.Runtime, filePath string, types string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(vm.ToValue(fmt.Errorf(`判断存在时出错：%s: %w`, filePath, err)))
	}

	orMatch := false
	andMatch := true

	for _, t := range types {
		switch t {
		case '*':
			orMatch = true
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

	return orMatch && andMatch
}

func exists(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	filePath := utils.MustBeString(call.Argument(0), vm)

	types := `*`
	if len(call.Arguments) > 2 {
		types = utils.MustBeString(call.Argument(1), vm)
	}

	return vm.ToValue(_exists(vm, filePath, types))
}

func fileExists(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	filePath := utils.MustBeString(call.Argument(0), vm)
	return vm.ToValue(_exists(vm, filePath, `f`))
}

func dirExists(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	dirPath := utils.MustBeString(call.Argument(0), vm)
	return vm.ToValue(_exists(vm, dirPath, `d`))
}

func mkDir(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	return _mkDir(call, vm, false)
}

func mkDirAll(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	return _mkDir(call, vm, true)
}

func _mkDir(call goja.FunctionCall, vm *goja.Runtime, all bool) goja.Value {
	path := utils.MustBeString(call.Argument(0), vm)

	perm := fs.FileMode(0755)
	if len(call.Arguments) >= 2 {
		var n int32
		if err := vm.ExportTo(call.Argument(1), &n); err != nil {
			panic(vm.ToValue(err))
		}
		perm = fs.FileMode(n)
	}

	fn := os.Mkdir
	if all {
		fn = os.MkdirAll
	}

	if err := fn(path, perm); err != nil {
		panic(vm.ToValue(err))
	}

	return nil
}
