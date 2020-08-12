package gw

import (
	"github.com/oceanho/gw/logger"
	"sync"
)

func Run(servers ...*HostServer) {
	var shutdown = sync.WaitGroup{}
	var started = sync.WaitGroup{}
	shutdown.Add(len(servers))
	started.Add(len(servers))
	for i := 0; i < len(servers); i++ {
		go func(s *HostServer) {
			go s.Serve()
			<-s.serverStartDone
			close(s.serverStartDone)
			started.Done()
			<-s.serverExitSignal
			close(s.serverExitSignal)
			logger.Info("Server: %s, Addr: %s exiting", s.options.Name, s.options.Addr)
			shutdown.Done()
		}(servers[i])
	}
	started.Wait()
	logger.Info("Successful, Waiting %d Servers exit", len(servers))
	shutdown.Wait()
	logger.Info("bye!")
}
