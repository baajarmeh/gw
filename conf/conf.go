package conf

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	json "github.com/json-iterator/go"
	"github.com/oceanho/gw/logger"
	"gopkg.in/yaml.v3"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type ApplicationConfigProvider interface {
	Provide(bcs BootConfig, out *ApplicationConfig) error
}

// ========================================================== //
//                                                            //
//    Define all of configuration Items for app Boot     //
//                                                            //                                                              ////                                                              ////
// ========================================================== //

// BootConfig represents a Application boot strap configuration object. It's likes linux boot.cnf
type BootConfig struct {
	Version                 string      `yaml:"version" toml:"version" json:"version"`
	ConfigProvider          string      `yaml:"configProvider" toml:"configProvider" json:"configProvider"`
	Configuration           interface{} `yaml:"configuration" toml:"configuration" json:"configuration"`
	configurationJsonString string
	locker                  sync.Mutex
}

func (bc *BootConfig) ParserTo(out interface{}) error {
	bc.locker.Lock()
	defer bc.locker.Unlock()
	if bc.configurationJsonString == "" {
		b, err := json.Marshal(bc.Configuration)
		if err != nil {
			return err
		}
		bc.configurationJsonString = string(b)
	}
	return json.Unmarshal([]byte(bc.configurationJsonString), out)
}

type GwConfigInfra struct {
	Addr     string            `yaml:"addr" toml:"addr" json:"addr"`
	AppID    string            `yaml:"appid" toml:"appid" json:"appid"`
	Secret   string            `yaml:"secret" toml:"secret" json:"secret"`
	Provider string            `yaml:"provider" toml:"provider" json:"provider"`
	Args     map[string]string `yaml:"args" toml:"args" json:"args"`
}

type LocalFile struct {
	Path string `yaml:"path" toml:"path" json:"path"`
}

func (bc BootConfig) String() string {
	b, _ := json.MarshalIndent(bc, "", "  ")
	return string(b)
}

//
// ======================================================= //
//                                                         //
//    Define all of configuration Items for Application    //
//                                                         //
// ======================================================= //
//
type ApplicationConfig struct {
	Server            Server      `yaml:"server" toml:"server" json:"server"`
	Backend           Backend     `yaml:"backend" toml:"backend" json:"backend"`
	Security          Security    `yaml:"security" toml:"security" json:"security"`
	Service           Service     `yaml:"service" toml:"service" json:"service"`
	Settings          Settings    `yaml:"settings" toml:"settings" json:"settings"`
	Custom            interface{} `yaml:"custom" toml:"custom" json:"custom"`
	locker            sync.Mutex
	customJson        string
	customMapState    int
	customJsonStrMaps map[string]string
}

// allow urls
type AllowUrl struct {
	Name string   `yaml:"name" toml:"name" json:"name"`
	Urls []string `yaml:"urls" toml:"urls" json:"urls"`
	IPs  []string `yaml:"ips" toml:"ips" json:"ips"`
}

//
// Server
//
type Server struct {
	Name       string `yaml:"name" toml:"name" json:"name"`
	ListenAddr string `yaml:"listenAddr" toml:"listenAddr" json:"listenAddr"`
}

//
// Backend Store Configuration
//
type Backend struct {
	Db    []Db    `yaml:"db" toml:"db" json:"db"`
	Cache []Cache `yaml:"cache" toml:"cache" json:"cache"`
}

// Backend of db
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

// Backend of cache
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

//
// Security
//
type Security struct {
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
			SidGenerator string `yaml:"sidGenerator" toml:"sidGenerator" json:"sidGenerator"`
		} `yaml:"session" toml:"session" json:"session"`
		Permission struct {
			DefaultStore struct {
				Name string `yaml:"name" toml:"name" json:"name"`
				Type string `yaml:"type" toml:"type" json:"type"`
			} `yaml:"defaultStore" toml:"defaultStore" json:"defaultStore"`
		} `yaml:"permission" toml:"permission" json:"permission"`
		Cookie struct {
			Key      string `yaml:"key" toml:"key" json:"key"`
			Path     string `yaml:"path" toml:"path" json:"path"`
			HttpOnly bool   `yaml:"httpOnly" toml:"httpOnly" json:"httpOnly"`
			MaxAge   int    `yaml:"maxAge" toml:"maxAge" json:"maxAge,string"`
			Secure   bool   `yaml:"secure" toml:"secure" json:"secure"`
			Domain   string `yaml:"domain" toml:"domain" json:"domain"`
		} `yaml:"cookie" toml:"cookie" json:"cookie"`
		AllowUrls []AllowUrl `yaml:"allowUrls" toml:"allowUrls" json:"allowUrls"`
	} `yaml:"auth" toml:"auth" json:"auth"`
	AuthServer struct {
		EnableAuthServe bool   `yaml:"enableAuthServe" toml:"enableAuthServe" json:"enableAuthServe"`
		Addr            string `yaml:"addr" toml:"addr" json:"addr"`
		LogIn           struct {
			Url       string   `yaml:"url" toml:"url" json:"url"`
			Methods   []string `yaml:"methods" toml:"methods" json:"methods"`
			AuthTypes []string `yaml:"authTypes" toml:"authTypes" json:"authTypes"`
		} `yaml:"login" toml:"login" json:"login"`
		LogOut struct {
			Url     string   `yaml:"url" toml:"url" json:"url"`
			Methods []string `yaml:"methods" toml:"methods" json:"methods"`
		} `yaml:"logout" toml:"logout" json:"logout"`
	} `yaml:"authServer" toml:"authServer" json:"authServer"`
	Limit struct {
		Pagination struct {
			MinPageSize int `yaml:"minPageSize" toml:"minPageSize" json:"minPageSize,string"`
			MaxPageSize int `yaml:"maxPageSize" toml:"maxPageSize" json:"maxPageSize,string"`
		} `yaml:"pagination" toml:"pagination" json:"pagination"`
	} `yaml:"limit" toml:"limit" json:"limit"`
}

//
// Service
//
type Service struct {
	Name    string `yaml:"name" toml:"name" json:"name"`
	Prefix  string `yaml:"prefix" toml:"prefix" json:"prefix"`
	Version string `yaml:"version" toml:"version" json:"version"`
	Remarks string `yaml:"remarks" toml:"remarks" json:"remarks"`
	PProf   struct {
		Enabled bool   `yaml:"enabled" toml:"enabled" json:"enabled"`
		Router  string `yaml:"router" toml:"router" json:"router"`
	} `yaml:"pprof" toml:"pprof" json:"pprof"`
	ServiceDiscovery struct {
		Enabled        bool `yaml:"enabled" toml:"enabled" json:"enabled"`
		RegistryCenter struct {
			Addr string `yaml:"addr" toml:"addr" json:"addr"`
		} `yaml:"registryCenter" toml:"registryCenter" json:"registryCenter"`
	} `yaml:"serviceDiscovery" toml:"serviceDiscovery" json:"serviceDiscovery"`
}

//
// Settings
//
type Settings struct {
	GwFramework struct {
		PrintRouterInfo struct {
			Disabled bool   `yaml:"disabled" toml:"disabled" json:"disabled"`
			Title    string `yaml:"title" toml:"title" json:"title"`
		} `yaml:"printRouterInfo" toml:"printRouterInfo" json:"printRouterInfo"`
	} `yaml:"gw" toml:"gw" json:"gw"`
	HeaderKey struct {
		RequestIDKey string `yaml:"requestIdKey" toml:"requestIdKey" json:"requestIdKey"`
	} `yaml:"headerKey" toml:"headerKey" json:"headerKey"`
	Monitor struct {
		SentryDNS string `yaml:"sentryDsn" toml:"sentryDsn" json:"sentryDsn"`
	} `yaml:"monitor" toml:"monitor" json:"monitor"`
	TimeoutControl struct {
		HTTP     int `yaml:"http" toml:"http" json:"http,string"`
		Redis    int `yaml:"redis" toml:"redis" json:"redis,string"`
		MongoDB  int `yaml:"mongo" toml:"mongo" json:"mongo,string"`
		Database int `yaml:"database" toml:"database" json:"database,string"`
	} `yaml:"timeoutControl" toml:"timeoutControl" json:"timeoutControl"`
	ExpirationTimeControl struct {
		Session int `yaml:"session" toml:"session" json:"session,string"`
	} `yaml:"expirationTimeControl" toml:"expirationTimeControl" json:"expirationTimeControl"`
}

func (cnf ApplicationConfig) String() string {
	b, _ := json.MarshalIndent(cnf, "", "  ")
	return string(b)
}

func genCustomMaps(prefix string, mapStore map[string]string, custom interface{}) {
	var mps, ok = custom.(map[interface{}]interface{})
	if ok {
		for k, v := range mps {
			processMaps(prefix, mapStore, k, v)
		}
	} else {
		var mps, ok = custom.(map[string]interface{})
		if ok {
			for k, v := range mps {
				processMaps(prefix, mapStore, k, v)
			}
		}
	}
}

func processMaps(prefix string, mapStore map[string]string, k interface{}, v interface{}) {
	k1, o := k.(string)
	if o {
		if prefix != "" {
			k1 = prefix + "." + k1
		}
		val := v.(interface{})
		b, e := json.Marshal(val)
		if e != nil {
			panic(fmt.Sprintf("genCustomMaps, err: %v", e))
		}
		mapStore[k1] = string(b)
		genCustomMaps(k1, mapStore, v)
	}
}

var (
	ErrNoPathSection = fmt.Errorf("no path section")
)

// ParseCustomPathTo parse the specifies path as out interface object.
//
// examples:
// Your custom section of app.yaml likes:
//
// comstom:
//   gwpro:
//     backend:
//       db: primary
//       cache: primary
//
// Now, your can use this API parse gwpro section to your custom go struct.
// It's looks as:
//
// type GwPro struct {
//   Backend struct {
//     Db string  `"json":"db" "toml":"db" "yaml":"db"`
//     Cache string  `"json":"cache" "toml":"cache" "yaml":"cache"`
//   } `"json":"backend" "toml":"backend" "yaml":"backend"`
// }
//
// var gwPro GwPro
//
// <Your ApplicationConfig>.ParseCustomPathTo("gwpro", &gwPro)
//
func (cnf *ApplicationConfig) ParseCustomPathTo(path string, out interface{}) error {
	str, ok := cnf.customJsonStrMaps[path]
	if !ok {
		return ErrNoPathSection
	}
	return json.Unmarshal([]byte(str), out)
}

func (cnf *ApplicationConfig) ParseCustomTo(out interface{}) error {
	return json.Unmarshal([]byte(cnf.customJson), out)
}

func (cnf *ApplicationConfig) Compile() {
	cnf.locker.Lock()
	defer cnf.locker.Unlock()
	if cnf.customMapState == 0 {
		cnf.customJsonStrMaps = make(map[string]string)
		genCustomMaps("", cnf.customJsonStrMaps, cnf.Custom)
		cnf.customMapState++
		b, err := json.Marshal(cnf.Custom)
		if err != nil {
			panic(fmt.Sprintf("compile application fail, err: %v", err))
		}
		cnf.customJson = string(b)
		cnf.customMapState++
	}
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

var (
	defaultBootConfigFileName = "config/boot.yaml"
	configProviders           map[string]ApplicationConfigProvider
	suffixParsers             map[string]func(b []byte, out interface{}) error
	TemplateParser            func(prepare interface{}, out *ApplicationConfig) error
)

func init() {

	configProviders = make(map[string]ApplicationConfigProvider)
	initialSuffixParsers()

	// config provider initialization
	localFileProvider := newLocalFileConfigProvider()
	RegisterProvider("lf", localFileProvider)
	RegisterProvider("local", localFileProvider)
	RegisterProvider("localfs", localFileProvider)
	RegisterProvider("localfile", localFileProvider)

	gwConfigProvider := newGWHttpConfSvrConfigProvider()
	RegisterProvider("gw-conf", gwConfigProvider)
	RegisterProvider("gw-infra", gwConfigProvider)
	RegisterProvider("gw-infra-conf", gwConfigProvider)

	TemplateParser = func(cnf interface{}, out *ApplicationConfig) error {
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

func initialSuffixParsers() {
	// initial
	suffixParsers = make(map[string]func(b []byte, out interface{}) error)
	// yaml
	suffixParsers[".yaml"] = func(b []byte, out interface{}) error {
		err := yaml.Unmarshal(b, out)
		return err
	}
	suffixParsers[".yml"] = suffixParsers[".yaml"]

	// toml
	suffixParsers[".toml"] = func(b []byte, out interface{}) error {
		err := toml.Unmarshal(b, out)
		return err
	}

	suffixParsers[".tml"] = suffixParsers[".toml"]

	// json
	suffixParsers[".json"] = func(b []byte, out interface{}) error {
		err := json.Unmarshal(b, out)
		return err
	}
}

func RegisterProvider(name string, provider ApplicationConfigProvider) {
	configProviders[name] = provider
}

func NewConfigWithBootConfig(bcs *BootConfig) *ApplicationConfig {
	cp, ok := configProviders[bcs.ConfigProvider]
	if !ok {
		panic(fmt.Sprintf("provider: %s are not support. you can added by gw.RegisterConfigProvider(...).", bcs.ConfigProvider))
	}
	cnf := &ApplicationConfig{}
	err := cp.Provide(*bcs, cnf)
	if err != nil {
		panic(fmt.Sprintf("config provider: %s call Provide(...) fail. err: %v", bcs.ConfigProvider, err))
	}
	return cnf
}

func NewBootConfigFromBytes(formatter string, bytes []byte) *BootConfig {
	logger.Debug("exec LoadBootConfigFromBytes(...), formatter: %s", formatter)
	formatter = strings.TrimLeft(formatter, ".")
	ext := fmt.Sprintf(".%s", formatter)
	p, o := suffixParsers[ext]
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

func NewBootConfigFromFile(filename string) *BootConfig {
	logger.Debug("exec LoadBootConfigFromFile(...) path: %s", filename)
	ext := filepath.Ext(filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("read boostrap conf file: %s, err: %v", filename, err))
	}
	return NewBootConfigFromBytes(ext, b)
}

// DefaultBootConfig returns a bcs from local file config/boot.yaml
func DefaultBootConfig() *BootConfig {
	return NewBootConfigFromFile(defaultBootConfigFileName)
}
