package jg_fs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/dop251/goja"
	loop "github.com/movsb/jg/runtime"
	"github.com/movsb/jg/utils"
)

func sha256sum(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	filePath := utils.MustBeString(call.Argument(0), vm)

	fp, err := os.Open(filePath)
	if err != nil {
		panic(vm.ToValue(fmt.Errorf(`failed to open file: %s: %w`, filePath, err)))
	}
	runtime.AddCleanup(fp, func(int) { fp.Close() }, 0)

	promise, resolve, reject := loop.CreatePromise(vm)

	go func() {
		hash := sha256.New()
		if _, err := io.Copy(hash, fp); err != nil {
			reject(vm.ToValue(err))
			return
		}
		sum := hash.Sum(nil)
		str := fmt.Sprintf(`%x`, sum)
		resolve(vm.ToValue(str))
	}()

	return vm.ToValue(promise)
}
