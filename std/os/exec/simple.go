package jg_exec

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dop251/goja"
	"github.com/movsb/jg/utils"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// $`a ${b} c ${d}`
// taggedCommand(stringsArray['a ',' c ', ”], b, c)
// 第一个参数是字符串数组（带raw属性），个数始终为插值数+1个。
func taggedCommand(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	obj := call.Argument(0).ToObject(vm)
	if obj == nil {
		panic(vm.ToValue(fmt.Errorf(`not tag`)))
	}
	if obj.Get(`raw`).Equals(goja.Undefined()) {
		panic(vm.ToValue(`exec.$ must be used as tagged template literal function.`))
	}

	// use placeholders ${N} to replace js ${} expressions
	var args []string
	interpolations := map[string]any{}
	values := call.Arguments[1:]
	for i := range len(values) {
		args = append(args, obj.Get(fmt.Sprint(i)).String())
		name := fmt.Sprintf(`${__%d}`, i)
		args = append(args, name)
		interpolations[fmt.Sprintf(`__%d`, i)] = values[i].Export()
	}
	args = append(args, obj.Get(fmt.Sprint(len(values))).String())
	cmdline := strings.Join(args, ``)

	cmd, err := shell(cmdline, interpolations)
	if err != nil {
		panic(vm.ToValue(err))
	}

	cmd2 := exec.Command(cmd.Command, cmd.Arguments...)
	cmd3 := &Command{underlying: cmd2}
	return vm.ToValue(cmd3)
}

func noBackground(stmt *syntax.Stmt) {
	if stmt.Background {
		panic(`cannot run in background`)
	}
}
func noNegated(stmt *syntax.Stmt) {
	if stmt.Negated {
		panic(`cannot test negated`)
	}
}
func noAssigns(call *syntax.CallExpr) {
	if len(call.Assigns) > 0 {
		panic(`no assigns allowed`)
	}
}

func shell(cmdline string, interpolations map[string]any) (_ *Cmd, outErr error) {
	defer utils.CatchAsError(&outErr)

	parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
	file, err := parser.Parse(strings.NewReader(cmdline), ``)
	if err != nil {
		return nil, fmt.Errorf(`failed to intermediate interpolation string: %s: %w`, cmdline, err)
	}
	if len(file.Stmts) > 1 {
		return nil, fmt.Errorf(`single command only`)
	}
	stmt0 := file.Stmts[0]
	noBackground(stmt0)
	noNegated(stmt0)

	cmd := &Cmd{}

	switch typed := stmt0.Cmd.(type) {
	default:
		panic(`unsupported command type`)
	case *syntax.CallExpr:
		utils.Must(call(cmd, stmt0, typed, interpolations))
	}

	return cmd, nil
}

func expandWord(word *syntax.Word, env _ReplacedInterpolationExpander) (string, error) {
	if lit := word.Lit(); lit != `` {
		return lit, nil
	}
	argName, err := expand.Literal(&expand.Config{
		Env:     env,
		NoUnset: true,
	}, word)
	if err != nil {
		return ``, err
	}
	value := env.ValueOf(argName)
	if value == nil {
		return ``, fmt.Errorf(`unknown argument: %s`, argName)
	}
	switch typed := value.(type) {
	default:
		return ``, fmt.Errorf(`unknown value type: %v`, value)
	case string:
		return typed, nil
	case int, bool, float64:
		return fmt.Sprint(typed), nil
	}
}

func call(cmd *Cmd, stmt *syntax.Stmt, expr *syntax.CallExpr, interpolations map[string]any) error {
	noAssigns(expr)

	env := _ReplacedInterpolationExpander{Known: interpolations}

	command, err := expandWord(expr.Args[0], env)
	if err != nil {
		return err
	}

	var args []string
	for _, arg := range expr.Args[1:] {
		expanded, err := expandWord(arg, env)
		if err != nil {
			return err
		}
		args = append(args, expanded)
	}

	cmd.Command = command
	cmd.Arguments = args

	if len(stmt.Redirs) > 0 {
		return fmt.Errorf(`redirect is not supported`)
	}

	// for _, r := range stmt.Redirs {
	// 	if err := redir(cmd, r, interpolations); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type _ReplacedInterpolationExpander struct {
	Known map[string]any
}

func (r _ReplacedInterpolationExpander) ValueOf(name string) any {
	return r.unwrapValues(name)
}

func (r _ReplacedInterpolationExpander) Get(name string) expand.Variable {
	if _, ok := r.Known[name]; !ok {
		return expand.Variable{}
	}
	return expand.Variable{
		Kind: expand.String,
		Set:  true,
		Str:  r.wrap(name),
	}
}
func (r _ReplacedInterpolationExpander) wrap(name string) string {
	return fmt.Sprintf(`__jg_%s_jg__`, name)
}

var reSplitWrapped = regexp.MustCompile(`(?U:__jg_.*_jg__)`)

func (r _ReplacedInterpolationExpander) unwrapValues(name string) any {
	// 如果替换后不为空，证明有其它字面字符，则参数值必须为字符串。
	n := 0
	empty := reSplitWrapped.ReplaceAllStringFunc(name, func(s string) string {
		n++
		return ``
	})
	if n == 0 {
		return name
	}

	// 单值时可以为任意类型。
	if n == 1 && empty == `` {
		var onlyValue any
		reSplitWrapped.ReplaceAllStringFunc(name, func(s string) string {
			s = strings.TrimPrefix(s, `__jg_`)
			s = strings.TrimSuffix(s, `_jg__`)
			v, ok := r.Known[s]
			if !ok {
				panic(fmt.Errorf(`no such value: %s`, s))
			}
			onlyValue = v
			return ``
		})
		return onlyValue
	}

	// 其它情况必须为基本类型。
	return reSplitWrapped.ReplaceAllStringFunc(name, func(s string) string {
		s = strings.TrimPrefix(s, `__jg_`)
		s = strings.TrimSuffix(s, `_jg__`)
		v, ok := r.Known[s]
		if !ok {
			panic(fmt.Errorf(`no such value: %s`, s))
		}
		switch typed := v.(type) {
		default:
			panic(fmt.Errorf(`unknown value type: %v`, v))
		case bool, int, string, float64:
			return fmt.Sprint(typed)
		}
	})
}

func (r _ReplacedInterpolationExpander) Each(predicate func(name string, vr expand.Variable) bool) {
	for k := range r.Known {
		if !predicate(k, expand.Variable{
			Kind: expand.String,
			Set:  true,
			Str:  r.wrap(k),
		}) {
			break
		}
	}
}

type Cmd struct {
	Command   string
	Arguments []string
	// StdinPath  string
	// StdinObj   io.Reader
	// StdoutPath string
	// StdoutObj  io.Writer
}

/*
func redir(cmd *Cmd, redir *syntax.Redirect, interpolations map[string]any) (outErr error) {
	defer utils.CatchAsError(&outErr)

	left := func() int {
		n := 1
		if redir.N != nil && redir.N.Value != `` {
			n = utils.Must1(strconv.Atoi(redir.N.Value))
		}
		return n
	}
	// env := _ReplacedInterpolationExpander{Known: interpolations}
	right := func() any {
		// 普通文件重定向。
		if lit := redir.Word.Lit(); lit != `` {
			return lit
		}
		if len(redir.Word.Parts) == 1 {
			if _, ok := redir.Word.Parts[0].(*syntax.ParamExp); ok {

			}
		}
		return utils.Must1(expand.Literal(nil, redir.Word))
	}

	switch redir.Op {
	default:
		return fmt.Errorf(`unsupported redirect: %v`, redir.Op)
	case syntax.RdrOut: // >
		switch left() {
		default:
			panic(`unknown left file`)
		case 1:
			switch typed := right().(type) {
			case string:
				cmd.StdoutPath = typed
			case io.Writer:
				cmd.StdoutObj = typed
			default:
				panic(`unknown stdout file`)
			}
		}
	}

	return nil
}
*/
