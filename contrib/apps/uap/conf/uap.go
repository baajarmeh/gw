package conf

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/conf"
)

type Uap struct {
	Auth       SectionConfig `json:"authManager" yaml:"authManager" toml:"authManager"`
	Session    SectionConfig `json:"sessionManager" yaml:"sessionManager" toml:"sessionManager"`
	Permission SectionConfig `json:"permissionManager" yaml:"permissionManager" toml:"permissionManager"`
	User       struct {
		SectionConfig
		Users []User `json:"initUsers" yaml:"initUsers" toml:"initUsers"`
	} `json:"userManager" yaml:"userManager" toml:"userManager"`
}

type User struct {
	User     string      `json:"user" yaml:"user" toml:"user"`
	Passport string      `json:"passport" yaml:"passport" toml:"passport"`
	Secret   string      `json:"secret" yaml:"secret" toml:"secret"`
	UserType gw.UserType `json:"type,string" yaml:"type" toml:"type"`
	TenantId uint64      `json:"tenantId,string" yaml:"tenantId" toml:"tenantId"`
	Desc     string      `json:"desc" yaml:"desc" toml:"desc"`
}

type SectionConfig struct {
	Backend struct {
		Name string `json:"name" yaml:"name" toml:"name"`
		//Prefix string `json:"prefix" yaml:"prefix" toml:"prefix"`
	} `json:"backend" yaml:"backend" toml:"backend"`
	Cache struct {
		Enable          bool   `json:"enable" yaml:"enable" toml:"enable"`
		Name            string `json:"name" yaml:"name" toml:"name"`
		Prefix          string `json:"prefix" yaml:"prefix" toml:"prefix"`
		ExpirationHours int    `json:"expirationHours,string" yaml:"expirationHours" toml:"expirationHours"`
	} `json:"cache" yaml:"cache" toml:"cache"`
}

func GetUAP(appConf *conf.ApplicationConfig) Uap {
	var cnf Uap
	if err := appConf.ParseCustomPathTo("gwcontrib.uap", &cnf); err != nil {
		panic(fmt.Sprintf("load gwcontrib.uap fail, err: %v", err))
	}
	return cnf
}
