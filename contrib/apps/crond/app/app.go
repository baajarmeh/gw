package app

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/crond/app/Dto"
	"github.com/oceanho/gw/contrib/apps/crond/app/Service"
	"github.com/oceanho/gw/logger"
	"sync"
	"time"
)

type App struct {
	info           gw.AppInfo
	registerFunc   func(router *gw.RouterGroup)
	useFunc        func(option *gw.ServerOption)
	onPrepareFunc  func(app *App, state *gw.ServerState)
	onStartFunc    func(app *App, state *gw.ServerState)
	onShutDownFunc func(app *App, state *gw.ServerState)

	onServerShutDown  chan bool
	queryJobsDuration time.Duration
	jobManager        Service.JobManager
	eventManager      gw.IEventManager
	queryJobFunc      func() *Dto.JobPager

	jobLocker       sync.Mutex
	jobQueue        map[string]Dto.Job
	waitEnQueueJobs chan Dto.Job

	jobProcessor     func()
	onFail           func(err error)
	onJobProcessFail func(job Dto.Job, err error)
}

type JobProcessFailEvent struct {
	job Dto.Job
}

func (q JobProcessFailEvent) MetaInfo() gw.EventMetaInfo {
	return gw.EventMetaInfo{
		Name:     fmt.Sprintf("job-%s", q.job.ID),
		Category: q.job.Name,
		Data:     q.job,
	}
}

func New() *App {
	var app = &App{
		info: struct {
			ID         uint64
			Name       string
			Router     string
			Descriptor string
		}{ID: 0, Name: "gw.cron", Router: "ocean/gw-cron", Descriptor: "Gw cron services."},
		registerFunc: func(router *gw.RouterGroup) {
		},
		onServerShutDown:  make(chan bool, 1),
		queryJobsDuration: time.Second * 10,
		waitEnQueueJobs:   make(chan Dto.Job, 128),
		jobManager:        Service.DefaultJobManager(),
		jobQueue:          make(map[string]Dto.Job),
		onFail: func(err error) {
			logger.Error("gw.cron fail, err: %v", err)
		},
		useFunc: func(option *gw.ServerOption) {
		},
		onPrepareFunc: func(app *App, state *gw.ServerState) {
			// register Job manager
			//state.DI().Register(Service.DefaultJobManager())
			app.eventManager = state.EventManager()
			app.onJobProcessFail = func(job Dto.Job, err error) {
				var e = app.eventManager.Publish(&JobProcessFailEvent{
					job: job,
				})
				if e != nil {
					logger.Error("publish event of <JobProcessFailEvent> fail, err: %v", e)
				}
			}
			app.queryJobFunc = func() *Dto.JobPager {
				from := time.Now()
				to := from.Add(app.queryJobsDuration * 2)
				expr := gw.PagerExpr{PageNumber: 1, PageSize: 16384}
				r, e := app.jobManager.Query(from, to, expr)
				if e != nil {
					app.onFail(e)
					return nil
				}
				return r
			}
			app.jobProcessor = func() {
				select {
				case j, o := <-app.waitEnQueueJobs:
					// chan closed.
					if !o {
						break
					}
					app.jobLocker.Lock()
					defer app.jobLocker.Unlock()
					if _, ok := app.jobQueue[j.ID]; !ok {
						go func(app *App, job Dto.Job) {
							// wait xxx, then exec job task
							d := job.ExecPlanAt.Sub(time.Now())
							if d < 0 {
								return
							}
							_ = <-time.After(d)
						}(app, j)
						app.jobQueue[j.ID] = j
					}
				}
			}
		},
		onStartFunc: func(app *App, state *gw.ServerState) {
			go func() {
				defer close(app.waitEnQueueJobs)
				defer close(app.onServerShutDown)
				select {
				case <-app.onServerShutDown:
					break
				case <-time.After(app.queryJobsDuration):
					jobPager := app.queryJobFunc()
					for _, job := range jobPager.Data {
						_job := job
						app.waitEnQueueJobs <- _job
					}
				}
			}()
			go app.jobProcessor()
		},
		onShutDownFunc: func(app *App, state *gw.ServerState) {
			app.onServerShutDown <- true
		},
	}
	return app
}

func (a *App) Meta() gw.AppInfo {
	return a.info
}

func (a *App) Register(router *gw.RouterGroup) {
	a.registerFunc(router)
}

func (a *App) Use(option *gw.ServerOption) {
	a.useFunc(option)
}

func (a *App) OnPrepare(state *gw.ServerState) {
	a.onPrepareFunc(a, state)
}

func (a *App) OnStart(state *gw.ServerState) {
	a.onStartFunc(a, state)
}

func (a *App) OnShutDown(state *gw.ServerState) {
	a.onShutDownFunc(a, state)
}
