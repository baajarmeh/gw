package gw

func IfNotNullPanic(obj interface{}) {
	if obj != nil {
		panic(obj)
	}
}

func IfNullPanic(obj interface{}) {
	if obj == nil {
		panic(obj)
	}
}

func IfNullThen(obj interface{}, handler func()) {
	if obj == nil {
		handler()
	}
}

func IfNotNullThen(obj interface{}, handler func()) {
	if obj != nil {
		handler()
	}
}
