package jg_archive

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/movsb/jg/runtime/loop"
	"github.com/movsb/jg/utils"
)

var Methods = map[string]any{
	`TarReader`: newTarReader,
	`ZipReader`: newZipReader,
}

type TarReader struct {
	underlying *tar.Reader
}

func newTarReader(call goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	input, ok := call.Argument(0).Export().(io.Reader)
	if !ok {
		panic(vm.ToValue(`第一个参数不可读`))
	}

	tr := &TarReader{
		underlying: tar.NewReader(input),
	}

	obj := vm.ToValue(tr).(*goja.Object)
	obj.SetPrototype(call.This.Prototype())

	return obj
}

func (tr *TarReader) ExtractTo(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	outputDir := utils.MustBeString(call.Argument(0), vm)
	stat, err := os.Stat(outputDir)
	if err != nil {
		panic(vm.ToValue(fmt.Errorf(`目录错误：%w`, err)))
	}
	if !stat.IsDir() {
		panic(vm.ToValue(fmt.Errorf(`不是目录：%s`, outputDir)))
	}

	promise, resolve, reject := loop.CreatePromise(vm)

	go func() {
		for {
			hdr, err := tr.underlying.Next()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				reject(vm.ToValue(fmt.Errorf(`读取时错误：%w`, err)))
				return
			}

			if !filepath.IsLocal(hdr.Name) {
				reject(vm.ToValue(fmt.Errorf(`错误的文件路径：%s`, hdr.Name)))
				return
			}

			switch hdr.Typeflag {
			default:
				reject(vm.ToValue(fmt.Errorf(`未知文件类型：%d`, hdr.Typeflag)))
				return
			case tar.TypeDir:
				if err := os.MkdirAll(filepath.Join(outputDir, hdr.Name), 0755); err != nil {
					reject(vm.ToValue(fmt.Errorf(`目录创建失败：%w`, err)))
					return
				}
			case tar.TypeReg:
				dir, _ := path.Split(hdr.Name)
				if err := os.MkdirAll(filepath.Join(outputDir, dir), 0755); err != nil {
					reject(vm.ToValue(fmt.Errorf(`目录创建失败：%w`, err)))
					return
				}

				outputFile := filepath.Join(outputDir, hdr.Name)
				fp, err := os.Create(outputFile)
				if err != nil {
					reject(vm.ToValue(fmt.Errorf(`文件创建失败：%w`, err)))
					return
				}

				if _, err := io.Copy(fp, tr.underlying); err != nil {
					fp.Close()
					os.Remove(outputFile)
					reject(vm.ToValue(fmt.Errorf(`文件创建失败：%w`, err)))
					return
				}

				if err := fp.Close(); err != nil {
					os.Remove(outputFile)
					reject(vm.ToValue(fmt.Errorf(`文件关闭失败：%w`, err)))
					return
				}
			}
		}

		resolve(goja.Undefined())
	}()

	return vm.ToValue(promise)
}
