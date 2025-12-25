package jr

import "github.com/movsb/jg/std"

type Option func(rt *Runtime)

func WithStd() Option {
	return func(rt *Runtime) {
		rt.Run(std.Init)
	}
}
