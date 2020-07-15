package conf

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

var f = `
service:
  name: "confsvr"
  prefix: "api/v1"
  version: "Version 1.0"
  remarks: "User Account Platform API Services."
  security:
    auth:
      disable: False
    allow-urls:
      - POST:${API_PREFIX}/ucp/account/login
common:
  backend:
    db:
    - name: primary
      driver: mysql
      addr: 127.0.0.1
      port: 3306
      user: root
      password: oceanho
      database: djdb
      ssl_mode: on
      ssl_cert: ap
      args:
        charset: utf8
        parseTime: True
    cache:
    - name: primary
      driver: redis
      addr: 127.0.0.1
      port: 6379
      type: redis
      auth:
        disable: True
        user: ocean
        database: 1
        password: oceanho
custom:
  user: OceanHo
  tags:
  - body
  - programmer
`

func TestConfig(t *testing.T) {
	var cnf = &Config{}
	yaml.Unmarshal([]byte(f), cnf)
	j, _ := json.Marshal(cnf.Custom)
	t.Logf("%s", string(j))
	for k, v := range cnf.Custom {
		t.Logf("%v=%v, v type is. %v", k, v, reflect.TypeOf(v))
	}
	t.Logf("%v", cnf)
}
