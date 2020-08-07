package gw

type Decorator struct {
	Catalog  string
	MetaData interface{}
	Before   DecoratorHandler
	After    DecoratorHandler
}

type DecoratorHandler func(ctx *Context) (friendlyMsg string, err error)

// helpers
func FilterDecorator(filter func(d Decorator) bool, decorators ...Decorator) []Decorator {
	var result []Decorator
	for _, dc := range decorators {
		if filter(dc) {
			result = append(result, dc)
		}
	}
	return result
}
