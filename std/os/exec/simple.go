package jg_exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/dop251/goja"
	"mvdan.cc/sh/v3/syntax"
)

func taggedCommand(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	obj := call.Argument(0).ToObject(vm)
	if obj == nil {
		panic(vm.ToValue(fmt.Errorf(`not tag`)))
	}
	if obj.Get(`raw`).Equals(goja.Undefined()) {
		panic(vm.ToValue(`exec.$ must be used as tagged template literal function.`))
	}

	var args []string

	// use placeholders ${N} to replace js ${} expressions
	interpolates := call.Arguments[1:]
	for i := range len(interpolates) {
		args = append(args, obj.Get(fmt.Sprint(i)).String())
		args = append(args, fmt.Sprintf(`${__%d}`, i))
	}
	args = append(args, obj.Get(fmt.Sprint(len(interpolates))).String())

	cmdline := strings.Join(args, ``)

	parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
	file, err := parser.Parse(strings.NewReader(cmdline), ``)
	if err != nil {
		panic(vm.ToValue(fmt.Errorf(`failed to intermediate interpolation string: %s`, cmdline)))
	}
	callExpr := file.Stmts[0].Cmd.(*syntax.CallExpr)
	if callExpr == nil {
		panic(vm.ToValue(fmt.Errorf(`only simple command supported now`)))
	}

	var recurse func(parts []syntax.WordPart)
	recurse = func(parts []syntax.WordPart) {
		for i, part := range parts {
			switch typed := part.(type) {
			default:
				panic(vm.ToValue(fmt.Sprintf(`unhandled word part: %T, %v`, part, part)))
			case *syntax.Lit:
				continue
			case *syntax.ParamExp:
				var n int
				if _, err := fmt.Sscanf(typed.Param.Value, `__%d`, &n); err != nil {
					panic(vm.ToValue(fmt.Errorf(`unknown expression: %s`, typed.Param.Value)))
				}
				parts[i] = &syntax.Lit{
					Value: interpolates[n].String(),
				}
			case *syntax.DblQuoted:
				recurse(typed.Parts)
			}
		}
	}

	for _, arg := range callExpr.Args {
		recurse(arg.Parts)
	}

	printer := syntax.NewPrinter()
	args = args[:0]
	for _, arg := range callExpr.Args {
		buf := bytes.NewBuffer(nil)
		if err := printer.Print(buf, arg); err != nil {
			panic(vm.ToValue(fmt.Errorf(`failed to print: %w`, err)))
		}
		args = append(args, buf.String())
	}

	myCmd := &Command{
		underlying: exec.Command(args[0], args[1:]...),
	}

	return vm.ToValue(myCmd).(*goja.Object)
}
