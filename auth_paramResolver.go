package gw

import "github.com/gin-gonic/gin"

type DefaultAuthParamResolverImpl struct {
}

func (d DefaultAuthParamResolverImpl) Resolve(ctx *gin.Context) AuthParameter {
	panic("implement me")
}

func DefaultAuthParamResolver(state *ServerState) IAuthParamResolver {
	var authParamResolver = &DefaultAuthParamResolverImpl{}
	return authParamResolver
}
