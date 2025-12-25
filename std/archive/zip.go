package jg_archive

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dop251/goja"
	"github.com/movsb/jg/runtime/loop"
	"github.com/movsb/jg/utils"
)

type ZipReader struct {
	underlying *zip.ReadCloser
}

// TODO 要读文件，看起来比较耗时？考虑转为 Promise。
func newZipReader(call goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	firstArg := call.Argument(0).Export()
	filePath, ok := firstArg.(string)
	if !ok {
		loop.Panic(vm, `第一个参数应为路径：%v`, firstArg)
	}

	zr, err := zip.OpenReader(filePath)
	if err != nil {
		loop.Panic(vm, `failed to zip: %w`, err)
	}

	mzr := &ZipReader{
		underlying: zr,
	}

	runtime.AddCleanup(mzr, func(int) { mzr.underlying.Close() }, 0)

	obj := vm.ToValue(mzr).(*goja.Object)
	obj.SetPrototype(call.This.Prototype())

	return obj
}

func (zr *ZipReader) ExtractTo(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	outputDir := utils.MustBeString(call.Argument(0), vm)
	stat, err := os.Stat(outputDir)
	if err != nil {
		loop.Panic(vm, `目录错误：%w`, err)
	}
	if !stat.IsDir() {
		loop.Panic(vm, `不是目录：%s`, outputDir)
	}

	// check if all files are in the same common prefix directory,
	// if true, remove that prefix. And check to see if all file names are valid.

	for _, file := range zr.underlying.File {
		if file.NonUTF8 {
			loop.Panic(vm, `only utf-8 file names are supported now: %x`, file.Name)
		}
		if !filepath.IsLocal(file.Name) {
			loop.Panic(vm, `invalid file name: %s`, file.Name)
		}
		file.Mode()
	}

	promise, resolve, reject := loop.CreatePromise(vm)

	go func() {
		for _, file := range zr.underlying.File {
			switch file.Mode().Type() {
			default:
				reject(vm.ToValue(fmt.Errorf(`unhandled file type`)))
				return
			case fs.ModeDir:
				path := filepath.Join(outputDir, file.Name)
				if err := os.MkdirAll(path, file.Mode().Perm()); err != nil {
					reject(vm.ToValue(fmt.Errorf(`failed to create dir: %s: %w`, path, err)))
					return
				}
			case 0:
				fr, err := file.Open()
				if err != nil {
					reject(vm.ToValue(fmt.Errorf(`failed to read zip: %s: %w`, file.Name, err)))
					return
				}
				defer fr.Close()

				// 目录一定比其内的文件先保存，所以不用再针对文件创建其外层目录。
				path := filepath.Join(outputDir, file.Name)
				fw, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode().Perm())
				if err != nil {
					reject(vm.ToValue(fmt.Errorf(`failed to create file: %s: %w`, file.Name, err)))
					return
				}

				digest := crc32.NewIEEE()
				tr := io.TeeReader(fr, digest)

				if _, err := io.Copy(fw, tr); err != nil {
					fw.Close()
					os.Remove(path)
					reject(vm.ToValue(fmt.Errorf(`failed to write file: %s: %w`, path, err)))
					return
				}

				if err := fw.Close(); err != nil {
					os.Remove(path)
					reject(vm.ToValue(fmt.Errorf(`failed to close file: %s: %w`, path, err)))
					return
				}

				if digest.Sum32() != file.CRC32 {
					os.Remove(path)
					reject(vm.ToValue(fmt.Errorf(`failed to checksum file: %s: crc32 mismatch`, path)))
					return
				}

				// error ignored.
				os.Chtimes(path, time.Time{}, file.Modified)
			}
		}

		resolve(goja.Undefined())
	}()

	return vm.ToValue(promise)
}
