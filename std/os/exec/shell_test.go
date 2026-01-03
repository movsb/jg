package jg_exec

import (
	"log"
	"slices"
	"strings"
	"testing"
)

func TestShell(t *testing.T) {
	for _, tc := range []struct {
		ID             string
		Cmdline        string
		Interpolations map[string]any
		Cmd            *Cmd
		Error          string
	}{
		{
			Cmdline: `${name}`,
			Cmd: &Cmd{
				Command: `echo`,
			},
			Interpolations: map[string]any{
				`name`: `echo`,
			},
		},
		{
			Cmdline: `echo 1 2 3`,
			Cmd: &Cmd{
				Command:   `echo`,
				Arguments: []string{`1`, `2`, `3`},
			},
		},
		{
			ID:      `echo_var`,
			Cmdline: `echo ${a}`,
			Error:   `a: unbound variable`,
		},
		{
			ID:      `echo_val`,
			Cmdline: `echo ${a} ${b} ${c} ${d}`,
			Interpolations: map[string]any{
				`a`: 1,
				`b`: `str`,
				`c`: true,
				`d`: 3.14,
			},
			Cmd: &Cmd{
				Command:   `echo`,
				Arguments: []string{`1`, `str`, `true`, `3.14`},
			},
		},
		{
			ID:      `echo_var_two`,
			Cmdline: `echo ${a}${b}`,
			Interpolations: map[string]any{
				`a`: `a`,
				`b`: 1,
			},
			Cmd: &Cmd{
				Command:   `echo`,
				Arguments: []string{`a1`},
			},
		},
		{
			ID:      `echo_var_interleave`,
			Cmdline: `echo 1${a}2${b}3`,
			Interpolations: map[string]any{
				`a`: `a`,
				`b`: 1,
			},
			Cmd: &Cmd{
				Command:   `echo`,
				Arguments: []string{`1a213`},
			},
		},
		{
			ID:      `echo_spaced`,
			Cmdline: `echo ${a}`,
			Interpolations: map[string]any{
				`a`: `spaced value`,
			},
			Cmd: &Cmd{
				Command:   `echo`,
				Arguments: []string{`spaced value`},
			},
		},
		{
			ID:      `echo_quoted`,
			Cmdline: `echo "lit" "${a} ${b}" ${c}`,
			Interpolations: map[string]any{
				`a`: `spaced value`,
				`b`: 2,
				`c`: `another value`,
			},
			Cmd: &Cmd{
				Command:   `echo`,
				Arguments: []string{`lit`, `spaced value 2`, `another value`},
			},
		},
		// {
		// 	ID:      `redir_rdrout_file`,
		// 	Cmdline: `echo 1 > f`,
		// 	Cmd: &Cmd{
		// 		Command:    `echo`,
		// 		Arguments:  []string{`1`},
		// 		StdoutPath: `f`,
		// 	},
		// },
		// {
		// 	ID:      `redir_rdrout_buffer`,
		// 	Cmdline: `echo 1 > ${buf}`,
		// 	Interpolations: map[string]any{
		// 		`buf`: bytes.NewBuffer(nil),
		// 	},
		// 	Cmd: &Cmd{
		// 		Command:   `echo`,
		// 		Arguments: []string{`1`},
		// 		StdoutObj: bytes.NewBuffer(nil),
		// 	},
		// },
	} {
		if tc.ID == `echo_quoted` {
			log.Println(`debug`)
		}
		cmd, err := shell(tc.Cmdline, tc.Interpolations)
		if err != nil {
			if tc.Error != `` && strings.Contains(err.Error(), tc.Error) {
				continue
			}
			t.Fatal(tc.Cmdline, err)
		}
		if cmd.Command != tc.Cmd.Command {
			t.Fatal(`command not equal`, tc.Cmd.Command, cmd.Command)
		}
		if !slices.Equal(cmd.Arguments, tc.Cmd.Arguments) {
			t.Fatal(`arguments not equal`, tc.Cmd.Arguments, cmd.Arguments)
		}
		// if cmd.StdoutPath != tc.Cmd.StdoutPath {
		// 	t.Fatal(`stdout not equal`, tc.Cmd.StdoutPath, cmd.StdoutPath)
		// }
	}
}
