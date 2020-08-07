package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/oceanho/gw/logger"
	"gopkg.in/yaml.v3"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type ConfigProvider interface {
	Name() string
	Provide(bcs BootConfig, out *Config) error
}

// ========================================================== //
//                                                            //
//    Define all of configuration Items for app Boot     //
//                                                            //                                                              ////                                                              ////
// ========================================================== //

// BootConfig represents a Application boot strap configuration object. It's likes linux boot.cnf
type BootConfig struct {
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

func (bsc BootConfig) String() string {
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
					MaxAge   int    `yaml:"maxAge" toml:"maxAge" json:"maxAge,string"`
					Secure   bool   `yaml:"secure" toml:"secure" json:"secure"`
					Domain   string `yaml:"domain" toml:"domain" json:"domain"`
				} `yaml:"cookie" toml:"cookie" json:"cookie"`
				Server struct {
					AuthServe bool   `yaml:"authServe" toml:"authServe" json:"authServe"`
					Addr      string `yaml:"addr" toml:"addr" json:"addr"`
					LogIn     struct {
						Url       string   `yaml:"url" toml:"url" json:"url"`
						Methods   []string `yaml:"methods" toml:"methods" json:"methods"`
						AuthTypes []string `yaml:"authTypes" toml:"authTypes" json:"authTypes"`
					} `yaml:"login" toml:"login" json:"login"`
					LogOut struct {
						Url     string   `yaml:"url" toml:"url" json:"url"`
						Methods []string `yaml:"methods" toml:"methods" json:"methods"`
					} `yaml:"logout" toml:"logout" json:"logout"`
				} `yaml:"server" toml:"server" json:"server"`
				AllowUrls []AllowUrl `yaml:"allowUrls" toml:"allowUrls" json:"allowUrls"`
			} `yaml:"auth" toml:"auth" json:"auth"`
			Limit struct {
				Pagination struct {
					MinPageSize int `yaml:"minPageSize" toml:"minPageSize" json:"minPageSize,string"`
					MaxPageSize int `yaml:"maxPageSize" toml:"maxPageSize" json:"maxPageSize,string"`
				} `yaml:"pagination" toml:"pagination" json:"pagination"`
			} `yaml:"limit" toml:"limit" json:"limit"`
		} `yaml:"security" toml:"security" json:"security"`
		Settings struct {
			GwFramework struct {
				PrintRouterInfo struct {
					Disabled bool   `yaml:"disabled" toml:"disabled" json:"disabled"`
					Title    string `yaml:"title" toml:"title" json:"title"`
				} `yaml:"printRouterInfo" toml:"printRouterInfo" json:"printRouterInfo"`
			} `yaml:"gw" toml:"gw" json:"gw"`
			HeaderKey struct {
				RequestIDKey string `yaml:"requestIdKey" toml:"requestIdKey" json:"requestIdKey"`
			} `yaml:"headerKey" toml:"headerKey" json:"headerKey"`
			TimeoutControl struct {
				HTTP     int `yaml:"http" toml:"http" json:"http,string"`
				Redis    int `yaml:"redis" toml:"redis" json:"redis,string"`
				MongoDB  int `yaml:"mongo" toml:"mongo" json:"mongo,string"`
				Database int `yaml:"database" toml:"database" json:"database,string"`
			} `yaml:"timeoutControl" toml:"timeoutControl" json:"timeoutControl"`
			ExpirationTimeControl struct {
				Session int `yaml:"session" toml:"session" json:"session,string"`
			} `yaml:"expirationTimeControl" toml:"expirationTimeControl" json:"expirationTimeControl"`
		} `yaml:"settings" toml:"settings" json:"settings"`
	} `yaml:"service" toml:"service" json:"service"`
	Common struct {
		Backend Backend `yaml:"backend" toml:"backend" json:"backend"`
	} `yaml:"common" toml:"common" json:"common"`
	Custom map[string]interface{} `yaml:"custom" toml:"custom" json:"custom"`
}

type AllowUrl struct {
	Name string   `yaml:"name" toml:"name" json:"name"`
	Urls []string `yaml:"urls" toml:"urls" json:"urls"`
	IPs  []string `yaml:"ips" toml:"ips" json:"ips"`
}

type Backend struct {
	Db    []Db    `yaml:"db" toml:"db" json:"db"`
	Cache []Cache `yaml:"cache" toml:"cache" json:"cache"`
}

type Db struct {
	Driver   string                 `yaml:"driver" toml:"driver" json:"driver"`
	Name     string                 `yaml:"name" toml:"name" json:"name"`
	Addr     string                 `yaml:"addr" toml:"addr" json:"addr"`
	Port     int                    `yaml:"port" toml:"port" json:"port,string"`
	User     string                 `yaml:"user" toml:"user" json:"user"`
	Password string                 `yaml:"password" toml:"password" json:"password"`
	Database string                 `yaml:"database" toml:"database" json:"database"`
	SSLMode  string                 `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"`
	SSLCert  string                 `yaml:"ssl_cert" toml:"ssl_cert" json:"ssl_cert"`
	Args     map[string]interface{} `yaml:"args" toml:"args" json:"args"`
}

type Cache struct {
	Driver           string                 `yaml:"driver" toml:"driver" json:"driver"`
	Name             string                 `yaml:"name" toml:"name" json:"name"`
	Addr             string                 `yaml:"addr" toml:"addr" json:"addr"`
	Port             int                    `yaml:"port" toml:"port" json:"port,string"`
	User             string                 `yaml:"user" toml:"user" json:"user"`
	Password         string                 `yaml:"password" toml:"password" json:"password"`
	DB               int                    `yaml:"db" toml:"db" json:"db,string"`
	SSLMode          string                 `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"`
	SSLCert          string                 `yaml:"ssl_cert" toml:"ssl_cert" json:"ssl_cert"`
	SSLCertFormatter string                 `yaml:"ssl_cert_fmt" toml:"ssl_cert_fmt" json:"ssl_cert_fmt"`
	Args             map[string]interface{} `yaml:"args" toml:"args" json:"args"`
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
	defaultBootConfigFileName = "config/boot.yaml"
	formatterDecoders         map[string]func(b []byte, out interface{}) error
	configProviders           map[string]ConfigProvider
	TemplateParser            func(prepare interface{}, out *Config) error
)

func init() {
	configProviders = make(map[string]ConfigProvider)
	formatterDecoders = make(map[string]func(b []byte, out interface{}) error)
	initialFormatterDecoders()

	// config provider initialization
	RegisterProvider(newLocalFileConfigProvider())
	RegisterProvider(newGWHttpConfSvrConfigProvider())

	TemplateParser = func(cnf interface{}, out *Config) error {
		b, e := json.Marshal(cnf)
		if e != nil {
			panic(fmt.Errorf("conf.TemplateParser(...) fail on json.Marshal(). err: %v", e))
		}
		s := string(b)
		tmpl, err := template.New("gw-config").Parse(s)
		if err != nil {
			panic(fmt.Sprintf("parse template fail. %v", err))
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, cnf)
		if err != nil {
			panic(fmt.Sprintf("execute template fail. %v", err))
		}
		return json.Unmarshal(buf.Bytes(), out)
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

func NewConfigByBootConfig(bcs *BootConfig) *Config {
	cp, ok := configProviders[bcs.AppCnf.Provider]
	if !ok {
		panic(fmt.Sprintf("provider: %s are not support. you can added by conf.RegisterProvider(...).", bcs.AppCnf.Provider))
	}
	cnf := &Config{}
	err := cp.Provide(*bcs, cnf)
	if err != nil {
		panic(fmt.Sprintf("config provider: %s call Provide(...) fail. err: %v", bcs.AppCnf.Provider, err))
	}
	return cnf
}

func LoadBootConfigFromBytes(formatter string, bytes []byte) *BootConfig {
	logger.Debug("exec LoadBootConfigFromBytes(...), formatter: %s", formatter)
	formatter = strings.TrimLeft(formatter, ".")
	ext := fmt.Sprintf(".%s", formatter)
	p, o := formatterDecoders[ext]
	if !o {
		panic(fmt.Sprintf("not supports bootstrap config suffix: %s.", ext))
	}
	out := &BootConfig{}
	err := p(bytes, out)
	if err != nil {
		panic(fmt.Sprintf("read boostrap conf, err: %v", err))
	}
	return out
}

func LoadBootConfigFromFile(filename string) *BootConfig {
	logger.Debug("exec LoadBootConfigFromFile(...) path: %s", filename)
	ext := filepath.Ext(filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("read boostrap conf file: %s, err: %v", filename, err))
	}
	return LoadBootConfigFromBytes(ext, b)
}

// DefaultBootConfig returns a bcs from local file config/boot.yaml
func DefaultBootConfig() *BootConfig {
	return LoadBootConfigFromFile(defaultBootConfigFileName)
}
