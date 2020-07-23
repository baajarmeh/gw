package conf

import (
	"testing"
)

var f = `
service:
  name: "confsvr"
  prefix: "/v1"
  version: "Version 1.0"
  remarks: "User Account Platform  Services."
  security:
    auth:
      disable: False
    allow-urls:
      - POST:${_PREFIX}/ucp/account/login
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

var bf = `
appconf:
  provider: localfs
  section: localfs
gwconf:
  addr: "https://configsvr.gw.com"
  appid: ""
  secret: ""
  type: plaintext
  provider: defaultHttpProvider
  args:
    salt: Salt$#NB
localfs:
  path: "config/app.yaml"
  type: plaintext
  formatter: yaml
`

func TestLoadBootStrapFromBytes_ShouldBe_OK(t *testing.T) {
	bsc := LoadBootStrapConfigFromBytes("yaml", []byte(bf))
	t.Logf("%v", bsc)
}

func TestLoadConfig_ShouldBe_OK(t *testing.T) {
	bsc := LoadBootStrapConfigFromBytes("yaml", []byte(bf))
	t.Logf("%v", bsc)
}
