package gw

func PanicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
