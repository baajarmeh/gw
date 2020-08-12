package conf

import "fmt"

type HttpProtoConfig struct {
	Addr          string
	Path          string
	Method        string
	Authorization struct {
		User     string
		Password string
	}
}

type HttpProtoConfigProvider struct {
}

func (h HttpProtoConfigProvider) Provide(bcs BootConfig, out *ApplicationConfig) error {
	var httpProtoConfig HttpProtoConfig
	if err := bcs.ParserTo(&httpProtoConfig); err != nil {
		panic(fmt.Sprintf("parser http proto config fail, err: %v", err))
	}

	panic("should be design and impl me please.")
}
