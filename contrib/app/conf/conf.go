package conf

type ConfigProvider interface {
	Read(out *Config) error
}

type ConfigReader interface {
	Auth(ak string, secret string)
}

// ======================================== //
//                                          //
//    Define all of configuration Items     //
//                                          //
// ======================================== //

type Config struct {
	Service struct {
		Name     string `yaml:"name" toml:"name" json:"name"`
		Prefix   string `yaml:"prefix" toml:"prefix" json:"prefix"`
		Version  string `yaml:"version" toml:"version" json:"version"`
		Remarks  string `yaml:"remarks" toml:"remarks" json:"remarks"`
		Security struct {
			Auth struct {
				Disable     bool     `yaml:"disable" toml:"disable" json:"disable"`
				AllowedUrls []string `yaml:"allow-urls" toml:"allow-urls" json:"allow-urls"`
			}
		} `yaml:"security" toml:"security" json:"security"`
	} `yaml:"service" toml:"service" json:"service"`
	Common struct {
		Backend *Backend `yaml:"backend" toml:"backend" json:"backend"`
	} `yaml:"common" toml:"common" json:"common"`
	Custom map[string]interface{} `yaml:"custom" toml:"custom" json:"custom"`
}

type Backend struct {
	Db    *[]Db    `yaml:"db"`
	Cache *[]Cache `yaml:"cache"`
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

// ============ End of configuration items ============= //
