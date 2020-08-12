package conf

import (
	"fmt"
	"io/ioutil"
	"path"
)

type GWHttpConfSvrConfigProvider struct {
}

func (c GWHttpConfSvrConfigProvider) Name() string {
	return "gwconf"
}

func (c GWHttpConfSvrConfigProvider) Provide(bcs BootConfig, out *ApplicationConfig) error {
	panic("implement me")
}

func newGWHttpConfSvrConfigProvider() GWHttpConfSvrConfigProvider {
	return GWHttpConfSvrConfigProvider{}
}

type LocalFileConfigProvider struct {
}

func (c LocalFileConfigProvider) Provide(bcs BootConfig, out *ApplicationConfig) error {
	var lf LocalFile
	err := bcs.ParserTo(&lf)
	if err != nil {
		panic(fmt.Sprintf("local file provider fail, err: %v", err))
	}
	if lf.Path == "" {
		panic(fmt.Sprintf("BootConfig missing configuration.Path section"))
	}
	suffix := path.Ext(lf.Path)
	provider, o := suffixParsers[suffix]
	if !o {
		return fmt.Errorf("not supports app config provider, suffix: %s", suffix)
	}
	b, err := ioutil.ReadFile(lf.Path)
	var outPrepare interface{}
	err = provider(b, &outPrepare)
	if err != nil {
		return fmt.Errorf("provider . %v", err)
	}
	return TemplateParser(outPrepare, out)
}

func newLocalFileConfigProvider() LocalFileConfigProvider {
	return LocalFileConfigProvider{}
}
