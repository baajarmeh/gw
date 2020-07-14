package conf

import "io"

type ConfigReader interface {
	io.Reader
	Auth(ak string, secret string)
}
