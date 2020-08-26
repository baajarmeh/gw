package generator

var tmpl = `
package {{ .Pkg }}

import (
	"github.com/oceanho/gw"
	_a "github.com/oceanho/gw/contrib/apps/uap"
)

const (
	IMyService1PkgName='github.com/oceanho/gw'
)

func GetIMyService1(ctx *gw.Context) _a.IMyService1() {
	return ctx.Resolve(IMyService1PkgName)
}
`
