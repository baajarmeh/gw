package gw

import "github.com/oceanho/gw/conf"

// RegisterConfigProvider defines a API that can be register your custom application config Provider.s
func RegisterConfigProvider(name string, provider conf.ApplicationConfigProvider) {
	conf.RegisterProvider(name, provider)
}
