package utils

import (
	"fmt"
	"reflect"

	"github.com/dop251/goja"
)

func MustBe[T any](value goja.Value, kind reflect.Kind, vm *goja.Runtime) T {
	if k := value.ExportType().Kind(); k != kind {
		panic(vm.ToValue(fmt.Sprintf(`参数类型不合法：%s vs %s`, kind, k)))
	}
	return value.Export().(T)
}

func MustBeString(value goja.Value, vm *goja.Runtime) string {
	return MustBe[string](value, reflect.String, vm)
}

func MustBeBoolean(value goja.Value, vm *goja.Runtime) bool {
	return MustBe[bool](value, reflect.Bool, vm)
}
