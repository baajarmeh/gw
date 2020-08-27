package gw

import (
	"fmt"
	"github.com/oceanho/gw/conf"
	"regexp"
)

type DefaultAuthParamCheckerImpl struct {
	config     *conf.ApplicationConfig
	validators map[string]*regexp.Regexp
}

var (
	ErrorAuthParamPassportFormatter   = fmt.Errorf("invalid auth parameter(passport)")
	ErrorAuthParamPasswordFormatter   = fmt.Errorf("invalid auth parameter(secret)")
	ErrorAuthParamVerifyCodeFormatter = fmt.Errorf("invalid auth parameter(verify code)")
)

func (d DefaultAuthParamCheckerImpl) Check(param AuthParameter) error {
	var pKey = d.config.Security.Auth.ParamKey
	if p, ok := d.validators[pKey.Passport]; ok {
		if !p.MatchString(param.Passport) {
			return ErrorAuthParamPassportFormatter
		}
	}
	if p, ok := d.validators[pKey.Secret]; ok {
		if !p.MatchString(param.Password) {
			return ErrorAuthParamPasswordFormatter
		}
	}
	if p, ok := d.validators[pKey.VerifyCode]; ok {
		if !p.MatchString(param.VerifyCode) {
			return ErrorAuthParamVerifyCodeFormatter
		}
	}
	return nil
}

func DefaultAuthParamChecker(state *ServerState) IAuthParamChecker {
	var authParamChecker = &DefaultAuthParamCheckerImpl{
		config:     state.ApplicationConfig(),
		validators: make(map[string]*regexp.Regexp),
	}
	p := authParamChecker.config.Security.Auth
	passportRegex, err := regexp.Compile(p.ParamPattern.Passport)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s for passport", p.ParamPattern.Passport))
	}
	secretRegex, err := regexp.Compile(p.ParamPattern.Secret)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s for secret", p.ParamPattern.Secret))
	}
	verifyCodeRegex, err := regexp.Compile(p.ParamPattern.VerifyCode)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s for verifyCode", p.ParamPattern.VerifyCode))
	}
	authParamChecker.validators[p.ParamKey.Passport] = passportRegex
	authParamChecker.validators[p.ParamKey.Secret] = secretRegex
	authParamChecker.validators[p.ParamKey.VerifyCode] = verifyCodeRegex
	return authParamChecker
}
