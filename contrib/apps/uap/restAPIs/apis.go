package restAPIs

import "github.com/oceanho/gw"

func DynamicAPIs() []gw.IDynamicRestAPI {
	var dynamicApis []gw.IDynamicRestAPI
	dynamicApis = append(dynamicApis, UserDynamicRestAPI())
	return dynamicApis
}
