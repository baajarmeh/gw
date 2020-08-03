package gwDb

import "fmt"

type ExprTree struct {
	Left     IFilter
	Right    IFilter
	Children []*ExprTree
}

type IFilter interface {
	Apply(column string) (expr string, params interface{})
}

type ISorter interface {
	Apply(column string) (expr string)
}

type SearchMode uint

const (
	Left SearchMode = iota
	Right
	FullSearch
)

type GreaterFilter struct {
	Columns []string
	Value   interface{}
}

type LessFilter struct {
	Columns []string
	Value   interface{}
}

type EqualFilter struct {
	Columns []string
	Value   interface{}
}

type LikeFilter struct {
	Columns    []string
	Keyword    string
	SearchMode SearchMode
}

func (f LikeFilter) Apply(column string) (expr string, params []interface{}) {
	params = make([]interface{}, 2)
	//params[0] = f.Left
	//params[1] = f.Right
	expr = fmt.Sprintf(" %s >= ? and %s <= ? ", column, column)
	return expr, params
}

type RangeFilter struct {
	Left  interface{}
	Right interface{}
}

func (f RangeFilter) Apply(column string) (expr string, params []interface{}) {
	params = make([]interface{}, 2)
	params[0] = f.Left
	params[1] = f.Right
	expr = fmt.Sprintf(" %s >= ? and %s <= ? ", column, column)
	return expr, params
}
