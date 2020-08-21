package conf

import "github.com/oceanho/gw"

type User struct {
	User     string      `json:"user" yaml:"user" toml:"user"`
	Passport string      `json:"passport" yaml:"passport" toml:"passport"`
	Secret   string      `json:"secret" yaml:"secret" toml:"secret"`
	UserType gw.UserType `json:"type,string" yaml:"type" toml:"type"`
	TenantId uint64      `json:"tenantId,string" yaml:"tenantId" toml:"tenantId"`
	Desc     string      `json:"desc" yaml:"desc" toml:"desc"`
}
