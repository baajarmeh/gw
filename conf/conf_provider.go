package conf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GWHttpConfSvrConfigProvider struct {
}

func (c GWHttpConfSvrConfigProvider) Name() string {
	return "gwconf"
}

func (c GWHttpConfSvrConfigProvider) Provide(bcs BootStrapConfig, out *Config) error {
	panic("implement me")
}

func newGWHttpConfSvrConfigProvider() GWHttpConfSvrConfigProvider {
	return GWHttpConfSvrConfigProvider{}
}

type LocalFileConfigProvider struct {
}

func (c LocalFileConfigProvider) Name() string {
	return "localfs"
}

func (c LocalFileConfigProvider) Provide(bcs BootStrapConfig, out *Config) error {
	formatter := bcs.LocalFS.Formatter
	if formatter == "" {
		formatter = filepath.Ext(bcs.LocalFS.Path)
	}
	formatter = strings.TrimLeft(formatter, ".")
	providerName := fmt.Sprintf(".%s", formatter)
	provider, o := formatterDecoders[providerName]
	if !o {
		return fmt.Errorf("not supports app config provider, suffix: %s", providerName)
	}
	b, err := ioutil.ReadFile(bcs.LocalFS.Path)
	err = provider(b, out)
	if err != nil {
		return fmt.Errorf("provider . %v", err)
	}
	return nil
}

func newLocalFileConfigProvider() LocalFileConfigProvider {
	return LocalFileConfigProvider{}
}
