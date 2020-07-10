package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App interface {
	Name() string
	Register(router *gin.Engine)
}

type Server struct {
	Addr                  string
	Name                  string
	Mode                  string
	Restart               string
	apps                  map[string]App
	locker                sync.Locker
	router                *gin.Engine
	StartBeforeHandler    func(server *Server) error
	ShutDownBeforeHandler func(server *Server) error
}

var (
	appDefaultAddr                  = ":8080"
	appDefaultName                  = "oceanho.app"
	appDefaultRestart               = "always"
	appDefaultMode                  = "release"
	appDefaultStartBeforeHandler    = func(server *Server) error { return nil }
	appDefaultShutDownBeforeHandler = func(server *Server) error { return nil }
)

func New() *Server {
	router := gin.Default()
	return &Server{
		Addr:                  appDefaultAddr,
		Name:                  appDefaultName,
		Restart:               appDefaultRestart,
		Mode:                  appDefaultMode,
		router:                router,
		StartBeforeHandler:    appDefaultStartBeforeHandler,
		ShutDownBeforeHandler: appDefaultShutDownBeforeHandler,
	}
}

func (server *Server) Register(app App) {
	server.locker.Lock()
	defer server.locker.Unlock()
	appName := app.Name()
	if _, ok := server.apps[appName]; !ok {
		server.apps[appName] = app
		app.Register(server.router)
	}
}

func (server *Server) Serve() {
	sigs := make(chan os.Signal, 1)
	handler := server.StartBeforeHandler
	if handler != nil {
		err := server.StartBeforeHandler(server)
		if err != nil {
			panic(fmt.Errorf("call app.StartBeforeHandler, %v", err))
		}
	}
	err := server.router.Run(server.Addr)
	if err != nil {
		panic(fmt.Errorf("call server.router.Run, %v", err))
	}
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	handler = server.ShutDownBeforeHandler
	if handler != nil {
		err := server.ShutDownBeforeHandler(server)
		if err != nil {
			fmt.Printf("call app.ShutDownBeforeHandler, %v", err)
		}
	}
}
