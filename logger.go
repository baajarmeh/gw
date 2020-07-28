package gw

import (
	"github.com/gin-gonic/gin"
	"os"
)

type Logger interface {
	Info(fmt interface{}, v ...interface{})
	Warn(fmt interface{}, v ...interface{})
	Debug(fmt interface{}, v ...interface{})
	Error(fmt interface{}, v ...interface{})
}

type DefaultImplLogger struct {
}

func (d DefaultImplLogger) Write(bytes []byte) (int, error) {
	return os.Stdout.Write(bytes)
}

func (d DefaultImplLogger) Info(fmt interface{}, v ...interface{}) {
	panic("implement me")
}

func (d DefaultImplLogger) Warn(fmt interface{}, v ...interface{}) {
	panic("implement me")
}

func (d DefaultImplLogger) Debug(fmt interface{}, v ...interface{}) {
	panic("implement me")
}

func (d DefaultImplLogger) Error(fmt interface{}, v ...interface{}) {
	panic("implement me")
}

type LoggerFactory struct {
}

func getLogger(ctx *gin.Context) Logger {
	// server := hostServer(ctx)
	return nil
}
