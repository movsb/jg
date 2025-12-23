package jg_path

import pathpkg "path"

type Path struct{}

func (*Path) Base(path string) string {
	return pathpkg.Base(path)
}

func init() {

}
