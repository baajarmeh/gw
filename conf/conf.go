package conf

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/oceanho/gw/logger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type ConfigProvider interface {
	Name() string
	Provide(bcs BootStrapConfig, out *Config) error
}

// ========================================================== //
//                                                            //
//    Define all of configuration Items for app BootStrap     //
//                                                            //                                                              ////                                                              ////
// ========================================================== //

type BootStrapConfig struct {
	AppCnf struct {
		Provider string `yaml:"provider" toml:"provider" json:"provider"`
		Section  string `yaml:"section" toml:"section" json:"section"`
	} `yaml:"appconf" toml:"appconf" json:"appconf"`
	GWConf struct {
		Addr     string            `yaml:"addr" toml:"addr" json:"addr"`
		AppID    string            `yaml:"appid" toml:"appid" json:"appid"`
		Secret   string            `yaml:"secret" toml:"secret" json:"secret"`
		Provider string            `yaml:"provider" toml:"provider" json:"provider"`
		Args     map[string]string `yaml:"args" toml:"args" json:"args"`
	} `yaml:"gwconf" toml:"gwconf" json:"gwconf"`
	LocalFS struct {
		Path      string `yaml:"path" toml:"path" json:"path"`
		Type      string `yaml:"type" toml:"type" json:"type"`
		Formatter string `yaml:"formatter" toml:"formatter" json:"formatter"`
	} `yaml:"localfs" toml:"localfs" json:"localfs"`
	Custom map[string]interface{} `yaml:"custom" toml:"custom" json:"custom"`
}

func (bsc BootStrapConfig) String() string {
	b, _ := json.MarshalIndent(bsc, "", "  ")
	return string(b)
}

// ======================================================= //
//                                                         //
//    Define all of configuration Items for Application    //
//                                                         //
// ======================================================= //

type Config struct {
	Service struct {
		Name     string `yaml:"name" toml:"name" json:"name"`
		Prefix   string `yaml:"prefix" toml:"prefix" json:"prefix"`
		Version  string `yaml:"version" toml:"version" json:"version"`
		Remarks  string `yaml:"remarks" toml:"remarks" json:"remarks"`
		Security struct {
			Crypto struct {
				Hash struct {
					Salt      string `yaml:"salt" toml:"salt" json:"salt"`
					Algorithm string `yaml:"alg" toml:"alg" json:"alg"`
				} `yaml:"hash" toml:"hash" json:"hash"`
				Protect struct {
					Secret    string `yaml:"secret" toml:"secret" json:"secret"`
					Algorithm string `yaml:"alg" toml:"alg" json:"alg"`
				} `yaml:"protect" toml:"protect" json:"protect"`
				Certificate struct {
					PrivateKey string `yaml:"privateKey" toml:"privateKey" json:"privateKey"`
					PublicKey  string `yaml:"publicKey" toml:"publicKey" json:"publicKey"`
					IsFile     bool   `yaml:"isFile" toml:"isFile" json:"isFile"`
				} `yaml:"cert" toml:"cert" json:"cert"`
			} `yaml:"crypto" toml:"crypto" json:"crypto"`
			Timeout struct {
				HTTP     int `yaml:"http" toml:"http" json:"http"`
				Redis    int `yaml:"redis" toml:"redis" json:"redis"`
				MongoDB  int `yaml:"mongo" toml:"mongo" json:"mongo"`
				Database int `yaml:"database" toml:"database" json:"database"`
			} `yaml:"timeout" toml:"timeout" json:"timeout"`
			Auth struct {
				TrustSidKey string `yaml:"trustSidKey" toml:"trustSidKey" json:"trustSidKey"`
				ParamKey    struct {
					Passport   string                 `yaml:"passport" toml:"passport" json:"passport"`
					Secret     string                 `yaml:"secret" toml:"secret" json:"secret"`
					VerifyCode string                 `yaml:"verifyCode" toml:"verifyCode" json:"verifyCode"`
					Custom     map[string]interface{} `yaml:"custom" toml:"custom" json:"custom"`
				} `yaml:"paramKey" toml:"paramKey" json:"paramKey"`
				ParamPattern struct {
					Passport   string `yaml:"passport" toml:"passport" json:"passport"`
					Secret     string `yaml:"secret" toml:"secret" json:"secret"`
					VerifyCode string `yaml:"verifyCode" toml:"verifyCode" json:"verifyCode"`
				} `yaml:"paramPattern" toml:"paramPattern" json:"paramPattern"`
				Session struct {
					DefaultStore struct {
						Name   string `yaml:"name" toml:"name" json:"name"`
						Prefix string `yaml:"prefix" toml:"prefix" json:"prefix"`
					} `yaml:"defaultStore" toml:"defaultStore" json:"defaultStore"`
				} `yaml:"session" toml:"session" json:"session"`
				Cookie struct {
					Key      string `yaml:"key" toml:"key" json:"key"`
					Path     string `yaml:"path" toml:"path" json:"path"`
					HttpOnly bool   `yaml:"httpOnly" toml:"httpOnly" json:"httpOnly"`
					MaxAge   int    `yaml:"maxAge" toml:"maxAge" json:"maxAge"`
					Secure   bool   `yaml:"secure" toml:"secure" json:"secure"`
					Domain   string `yaml:"domain" toml:"domain" json:"domain"`
				} `yaml:"cookie" toml:"cookie" json:"cookie"`
				AuthServer string     `yaml:"authServer" toml:"authServer" json:"authServer"`
				LoginUrl   string     `yaml:"loginUrl" toml:"loginUrl" json:"loginUrl"`
				LogoutUrl  string     `yaml:"logoutUrl" toml:"logoutUrl" json:"logoutUrl"`
				Disable    bool       `yaml:"disable" toml:"disable" json:"disable"`
				AllowUrls  []AllowUrl `yaml:"allowUrls" toml:"allowUrls" json:"allowUrls"`
			} `yaml:"auth" toml:"auth" json:"auth"`
			QueryLimit struct {
				MinPageSize int `yaml:"minPageSize" toml:"minPageSize" json:"minPageSize"`
				MaxPageSize int `yaml:"maxPageSize" toml:"maxPageSize" json:"maxPageSize"`
			} `yaml:"queryLimit" toml:"queryLimit" json:"queryLimit"`
		} `yaml:"security" toml:"security" json:"security"`
	} `yaml:"service" toml:"service" json:"service"`
	Common struct {
		Backend *Backend `yaml:"backend" toml:"backend" json:"backend"`
	} `yaml:"common" toml:"common" json:"common"`
	Custom map[string]interface{} `yaml:"custom" toml:"custom" json:"custom"`
}

type AllowUrl struct {
	Name string   `yaml:"name" toml:"name" json:"name"`
	Urls []string `yaml:"urls" toml:"urls" json:"urls"`
	IPs  []string `yaml:"ips" toml:"ips" json:"ips"`
}

type Backend struct {
	Db    []Db    `yaml:"db"`
	Cache []Cache `yaml:"cache"`
}

type Db struct {
	Driver   string            `yaml:"driver" toml:"driver" json:"driver"`
	Name     string            `yaml:"name" toml:"name" json:"name"`
	Addr     string            `yaml:"addr" toml:"addr" json:"addr"`
	Port     int               `yaml:"port" toml:"port" json:"port"`
	User     string            `yaml:"user" toml:"user" json:"user"`
	Password string            `yaml:"password" toml:"password" json:"password"`
	Database string            `yaml:"database" toml:"database" json:"database"`
	SSLMode  string            `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"`
	SSLCert  string            `yaml:"ssl_cert" toml:"ssl_cert" json:"ssl_cert"`
	Args     map[string]string `yaml:"args" toml:"args" json:"args"`
}

type Cache struct {
	Driver           string            `yaml:"driver" toml:"driver" json:"driver"`
	Name             string            `yaml:"name" toml:"name" json:"name"`
	Addr             string            `yaml:"addr" toml:"addr" json:"addr"`
	Port             int               `yaml:"port" toml:"port" json:"port"`
	User             string            `yaml:"user" toml:"user" json:"user"`
	Password         string            `yaml:"password" toml:"password" json:"password"`
	DB               int               `yaml:"db" toml:"db" json:"db"`
	SSLMode          string            `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"`
	SSLCert          string            `yaml:"ssl_cert" toml:"ssl_cert" json:"ssl_cert"`
	SSLCertFormatter string            `yaml:"ssl_cert_fmt" toml:"ssl_cert_fmt" json:"ssl_cert_fmt"`
	Args             map[string]string `yaml:"args" toml:"args" json:"args"`
}

func (cnf Config) String() string {
	b, _ := json.MarshalIndent(cnf, "", "  ")
	return string(b)
}

// ============ End of configuration items ============= //

// Extension defines.
type AllowUrlPattern struct {
	Method          string
	Pattern         string
	IsRegex         bool
	IPLimit         string
	regMatchPattern *regexp.Regexp
}

var defaultBSCLocalFileData = `
appconf:
  provider: localfs
  section: localfs
localfs:
  path: "config/app.yaml"
  type: plaintext
  formatter: yaml
`
var (
	defaultBootStrapConfigFileName = "config/boot.yaml"
	formatterDecoders              map[string]func(b []byte, out interface{}) error
	configProviders                map[string]ConfigProvider
	configVarModifier              func(cnf *Config) *Config
)

func init() {
	configProviders = make(map[string]ConfigProvider)
	formatterDecoders = make(map[string]func(b []byte, out interface{}) error)
	initialFormatterDecoders()

	// config provider initialization
	RegisterProvider(newLocalFileConfigProvider())
	RegisterProvider(newGWHttpConfSvrConfigProvider())

	configVarModifier = func(cnf *Config) *Config {
		b, e := json.Marshal(cnf)
		if e != nil {
			panic(fmt.Errorf("configVarModifier fail on json.Marshal(). err: %v", e))
		}
		s := string(b)
		prefix := cnf.Service.Prefix
		s = strings.Replace(s, "$PREFIX", prefix, -1)
		s = strings.Replace(s, "${PREFIX}", prefix, -1)
		cnfNew := &Config{}
		if json.Unmarshal([]byte(s), cnfNew) != nil {
			panic(fmt.Errorf("configVarModifier fail on json.Unmarshal(). err: %v", e))
		}
		return cnfNew
	}
}

func initialFormatterDecoders() {
	formatterDecoders[".yaml"] = func(b []byte, out interface{}) error {
		err := yaml.Unmarshal(b, out)
		return err
	}
	formatterDecoders[".yml"] = formatterDecoders[".yaml"]
	formatterDecoders[".toml"] = func(b []byte, out interface{}) error {
		err := toml.Unmarshal(b, out)
		return err
	}
	formatterDecoders[".tml"] = formatterDecoders[".toml"]
	formatterDecoders[".json"] = func(b []byte, out interface{}) error {
		err := json.Unmarshal(b, out)
		return err
	}
}

func RegisterProvider(provider ConfigProvider) {
	configProviders[provider.Name()] = provider
}

func DefaultFromGWConfSvr(cnf BootStrapConfig) *Config {
	panic("impl me.")
}

func NewConfigFromLocalFile(filename string) *Config {
	bcs := LoadBootStrapConfigFromBytes("yaml", []byte(defaultBSCLocalFileData))
	return NewConfigByBootStrapConfig(bcs)
}

func NewConfigByBootStrapConfig(bcs *BootStrapConfig) *Config {
	cp, ok := configProviders[bcs.AppCnf.Provider]
	if !ok {
		panic(fmt.Sprintf("provider: %s are not support. you can added by conf.RegisterProvider(...).", bcs.AppCnf.Provider))
	}
	cnf := &Config{}
	err := cp.Provide(*bcs, cnf)
	if err != nil {
		panic(fmt.Sprintf("config provider: %s call Provide(...) fail. err: %v", bcs.AppCnf.Provider, err))
	}
	return configVarModifier(cnf)
}

func LoadBootStrapConfigFromBytes(formatter string, bytes []byte) *BootStrapConfig {
	logger.Debug("exec LoadBootStrapConfigFromBytes(...), formatter: %s", formatter)
	formatter = strings.TrimLeft(formatter, ".")
	ext := fmt.Sprintf(".%s", formatter)
	p, o := formatterDecoders[ext]
	if !o {
		panic(fmt.Sprintf("not supports bootstrap config suffix: %s.", ext))
	}
	out := &BootStrapConfig{}
	err := p(bytes, out)
	if err != nil {
		panic(fmt.Sprintf("read boostrap conf, err: %v", err))
	}
	return out
}

func LoadBootStrapConfigFromFile(filename string) *BootStrapConfig {
	logger.Debug("exec LoadBootStrapConfigFromFile(...) path: %s", filename)
	ext := filepath.Ext(filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("read boostrap conf file: %s, err: %v", filename, err))
	}
	return LoadBootStrapConfigFromBytes(ext, b)
}

// DefaultBootStrapConfig returns a bcs from local file config/boot.yaml
func DefaultBootStrapConfig() *BootStrapConfig {
	return LoadBootStrapConfigFromFile(defaultBootStrapConfigFileName)
}
