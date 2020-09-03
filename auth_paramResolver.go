package gw

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserPasswordAuthParamResolver struct {
}

func (u UserPasswordAuthParamResolver) Resolve(c *gin.Context) AuthParameter {
	s := getHostServer(c)
	var tenantIdStr = ""
	var param AuthParameter
	paramKey := s.conf.Security.Auth.ParamKey
	// 1. User/Password
	param.CredType = UserPasswordAuth
	param.Passport, _ = c.GetPostForm(paramKey.Passport)
	param.Password, _ = c.GetPostForm(paramKey.Secret)
	param.VerifyCode, _ = c.GetPostForm(paramKey.VerifyCode)
	tenantIdStr, _ = c.GetPostForm(paramKey.TenantID)
	if check(param) {
		param.TenantID, _ = strconv.ParseUint(tenantIdStr, 10, 64)
		return param
	}

	// json
	if c.ContentType() == "application/json" {
		// json decode
		cred := gin.H{
			paramKey.Passport:   "",
			paramKey.Secret:     "",
			paramKey.VerifyCode: "",
			paramKey.TenantID:   "",
		}
		err := c.Bind(&cred)
		if err != nil {
			return param
		}
		param.Passport = cred[paramKey.Passport].(string)
		param.Password = cred[paramKey.Secret].(string)
		param.VerifyCode = cred[paramKey.VerifyCode].(string)
		tenantIdStr = cred[paramKey.TenantID].(string)
		if check(param) {
			param.TenantID, _ = strconv.ParseUint(tenantIdStr, 10, 64)
			return param
		}
	}
	return param
}

type BasicAuthParamResolver struct {
}

func (u BasicAuthParamResolver) Resolve(c *gin.Context) AuthParameter {
	s := getHostServer(c)
	var tenantIdStr = ""
	paramKey := s.conf.Security.Auth.ParamKey
	var param AuthParameter
	// 2. Basic auth
	var ok = false
	param.CredType = BasicAuth
	param.Passport, param.Password, ok = c.Request.BasicAuth()
	if param.VerifyCode == "" {
		param.VerifyCode = c.Query(paramKey.VerifyCode)
	}
	if ok && check(param) {
		param.TenantID, _ = strconv.ParseUint(tenantIdStr, 10, 64)
		return param
	}
	return param
}

type AksAuthParamResolver struct {
}

func (u AksAuthParamResolver) Resolve(c *gin.Context) AuthParameter {
	s := getHostServer(c)
	var param AuthParameter
	paramKey := s.conf.Security.Auth.ParamKey
	param.CredType = AksAuth
	var tenantIdStr = c.GetHeader(paramKey.TenantID)
	param.Passport = c.GetHeader(paramKey.Passport)
	param.Password = c.GetHeader(paramKey.Secret)
	param.VerifyCode = c.GetHeader(paramKey.VerifyCode)
	if check(param) {
		param.TenantID, _ = strconv.ParseUint(tenantIdStr, 10, 64)
		return param
	}
	return param
}

func check(param AuthParameter) bool {
	return param.Passport != "" && param.Password != ""
}

func DefaultAuthParamResolver() []IAuthParamResolver {
	var resolvers = make([]IAuthParamResolver, 3)
	resolvers[0] = &UserPasswordAuthParamResolver{}
	resolvers[1] = &BasicAuthParamResolver{}
	resolvers[2] = &AksAuthParamResolver{}
	return resolvers
}
