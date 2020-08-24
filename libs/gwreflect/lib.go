package gwreflect

import (
	"fmt"
	"reflect"
)

func GetPkgFullName(typer reflect.Type) string {
	if typer.Kind() == reflect.Ptr {
		return fmt.Sprintf("%s.%s", typer.Elem().PkgPath(), typer.Elem().Name())
	}
	return fmt.Sprintf("%s.%s", typer.PkgPath(), typer.Name())
}
