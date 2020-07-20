package conf

type GWHttpSvrConfig struct {
	AppID  string
	Secret string
}

func (c GWHttpSvrConfig) Provide(out *Config) error {
	panic("implement me")
}

type LocalFSConfig struct {
	Path string
}

func (c LocalFSConfig) Provide(out *Config) error {
	panic("implement me")
}
