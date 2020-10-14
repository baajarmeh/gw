package gw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oceanho/gw/conf"
	"github.com/oceanho/gw/logger"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

const (
	gwDbContextKey = "gw-context"
)

var (
	ErrorParamNotFound = fmt.Errorf("param not found")
)

// Handler defines a http handler for gw framework.
type Handler func(ctx *Context)

// Context represents a gw Context object, it's extension from gin.Context.
type Context struct {
	*gin.Context
	requestId  string
	user       User
	store      IStore
	startTime  time.Time
	logger     Logger
	queries    map[string][]string
	params     map[string]interface{}
	bindModels map[string]interface{}
	server     *HostServer
}

func (c *Context) RequestId() string {
	return c.requestId
}

func (c *Context) User() User {
	return c.user
}

func (c *Context) Store() IStore {
	return c.store
}

func (c *Context) Server() *HostServer {
	return c.server
}

func (c *Context) ResolveByTyper(typer reflect.Type) (err error, result interface{}) {
	err, result = c.server.DIProvider.ResolveByTyperWithState(c, typer)
	if err != nil {
		c.JSON400Msg(400, "resolve object fail")
	}
	return
}

func (c *Context) ResolveByObjectTyper(object interface{}) error {
	err := c.server.DIProvider.ResolveByObjectTyper(object)
	if err != nil {
		c.JSON400Msg(400, "resolve object fail")
	}
	return err
}

// StartTime returns the Context start *time.Time
func (c *Context) StartTime() *time.Time {
	return &c.startTime
}

// Config returns a snapshot of the current Context's conf.Config object.
func (c *Context) Config() *conf.ApplicationConfig {
	return config(c.Context)
}

// Query returns a string from queries.
func (c *Context) Query(key string) string {
	val := c.QueryArray(key)
	if len(val) > 0 {
		return val[0]
	}
	return ""
}

// QueryArray returns a string array from queries.
func (c *Context) QueryArray(key string) []string {
	val := c.Context.QueryArray(key)
	if len(val) == 0 {
		val = c.queries[key]
	}
	return val
}

// MustGetUint64IDFromParam returns a string from c.Params
func (c *Context) MustGetUint64IDFromParam(out *uint64) (err error) {
	var _out string
	if err := c.MustGetIdStrFromParam(&_out); err != nil {
		c.JSON400Msg(400, err)
		return err
	}
	var _outUint, err1 = strconv.ParseUint(_out, 10, 64)
	if _outUint < 1 {
		c.JSON400Msg(400, ErrorInvalidParamID)
		return ErrorInvalidParamID
	}
	*out = _outUint
	return err1
}

// MustGetIdStrFromParam returns a string from c.Params
func (c *Context) MustGetIdStrFromParam(out *string) error {
	return c.MustParam("id", out)
}

// MustParam returns a string from c.Params
func (c *Context) MustParam(key string, out *string) error {
	*out = c.Param(key)
	if *out == "" {
		c.JSON400Msg(400, fmt.Sprintf("missing parameter(%s)", key))
		return ErrorParamNotFound
	}
	return nil
}

// Bind define a Api that can be bind data to out object by gin.Context's Bind(...) APIs.
// It's auto response 400, invalid request parameter to client if bind fail.
// returns a error message for c.Bind(...).
func (c *Context) Bind(out interface{}) error {
	if err := c.Context.Bind(out); err != nil {
		c.JSON400Msg(400, "invalid request parameters")
		return err
	}
	return nil
}

// BindQuery define a Api that can be bind data to out object by gin.Context's Bind(...) APIs.
// It's auto response 400, invalid request parameter to client if bind fail.
// returns a error message for c.BindQuery(...).
func (c *Context) BindQuery(out interface{}) error {
	if err := c.Context.BindQuery(out); err != nil {
		c.JSON400Msg(400, "invalid request parameters")
		return err
	}
	return nil
}

func (c *Context) AppConfig() *conf.ApplicationConfig {
	return c.server.conf
}

// handle code APIs.
func handle(c *gin.Context) {
	var router, ok = c.MustGet(gwRouterInfoKey).(RouterInfo)
	if !ok {
		logger.Error("invalid handler, can not be got RouterInfo.")
		return
	}
	var status int
	var err error
	var shouldStop bool = false
	var payload interface{}
	// action before Decorators
	var s = getHostServer(c)
	var requestID = getRequestId(s, c)
	var ctx = makeCtx(c, requestID)
	for _, d := range router.beforeDecorators {
		status, err, payload = d.Before(ctx)
		if err != nil || status != 0 {
			shouldStop = true
			break
		}
	}
	if shouldStop {
		if payload == "" {
			payload = "caller decorator fail."
		}
		status, body := parseErrToRespBody(c, s, requestID, payload, err)
		c.JSON(status, body)
		return
	}

	// process Action handler.
	router.Handler(ctx)

	// action after Decorators
	l := len(router.afterDecorators)
	if l < 1 {
		return
	}
	shouldStop = false
	for i := l - 1; i >= 0; i-- {
		status, err, payload = router.afterDecorators[i].After(ctx)
		if err != nil || status != 0 {
			shouldStop = true
			break
		}
	}
	if shouldStop {
		if payload == "" {
			payload = "caller decorator fail."
		}
		status, body := parseErrToRespBody(c, s, requestID, payload, err)
		c.JSON(status, body)
	}
}

func parseErrToRespBody(c *gin.Context, s *HostServer, requestID string, msgBody interface{}, err error) (int, interface{}) {
	var status = http.StatusBadRequest
	if err == ErrPermissionDenied {
		status = http.StatusForbidden
	} else if err == ErrInternalServerError {
		status = http.StatusInternalServerError
	} else if err == ErrUnauthorized {
		status = http.StatusUnauthorized
	}
	return status, s.RespBodyBuildFunc(c, status, requestID, err.Error(), msgBody)
}

func makeCtx(c *gin.Context, requestID string) *Context {
	s := getHostServer(c)
	user := getUser(c)
	serverState := s.State()
	ctx := &Context{
		Context:    c,
		user:       user,
		server:     s,
		requestId:  requestID,
		startTime:  time.Now(),
		logger:     getLogger(c),
		queries:    make(map[string][]string),
		params:     make(map[string]interface{}),
		bindModels: make(map[string]interface{}),
	}
	var dbSetups []StoreDbSetupHandler
	dbSetups = append(dbSetups, s.storeDbSetupHandler)
	var cacheSetups []StoreCacheSetupHandler
	cacheSetups = append(cacheSetups, s.storeCacheSetupHandler)
	store := &backendWrapper{
		user:                    user,
		ctx:                     ctx,
		storeDbSetupHandlers:    dbSetups,
		storeCacheSetupHandlers: cacheSetups,
		store:                   serverState.Store(),
	}
	ctx.store = store
	return ctx
}

func splitDecorators(decorators ...Decorator) (before, after []Decorator) {
	for _, d := range decorators {
		if d.After != nil {
			after = append(after, d)
		}
		if d.Before != nil {
			before = append(before, d)
		}
	}
	return before, after
}
