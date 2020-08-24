package gw

import (
	"github.com/go-playground/assert/v2"
	"gorm.io/gorm"
	"testing"
)

type Dto1 struct {
}

type MyUser struct {
}

const (
	IMyServiceName string = "github.com/oceanho/gw.IMyService"
	MyService2Name string = "github.com/oceanho/gw.MyService2"
)

type IMyService interface {
	Create(dto Dto1) error
}

type MyService struct {
	User *gorm.DB
}

func (myService MyService) Create(dto Dto1) error {
	return nil
}

//
// Gw DI framework dependOns
func (MyService) New(store IStore) IMyService {
	var myService MyService
	myService.User = store.GetDbStore().Model(MyUser{})
	return myService
}

func (MyService) Destroy() {

}

type MyService2 struct {
	User       *gorm.DB
	MyService1 IMyService
}

func (myService MyService2) Create(dto Dto1) error {
	return nil
}

//
// Gw DI framework dependOns
func (MyService2) New(serverState ServerState, store IStore, service IMyService) MyService2 {
	var myService = MyService2{}
	myService.User = store.GetDbStore().Model(MyUser{})
	myService.MyService1 = service
	return myService
}

func (MyService2) Destroy() {

}

func TestDefaultDIProviderImpl(t *testing.T) {
	var server = DefaultServer()
	var state = NewServerState(server)
	var di = DefaultDIProvider(state)
	var myS1Impl MyService
	var myS2Impl MyService2
	di.Register(myS1Impl)
	di.Register(myS2Impl)
	var myS1, ok = di.Resolve(IMyServiceName).(IMyService)
	assert.IsEqual(ok, true)
	var dto Dto1
	_ = myS1.Create(dto)

	// Resolve can not return a pointer Object.
	var mys2 = di.Resolve(MyService2Name)
	// FIXME(OceanHo) ref #OC-20200824.001
	//var myS2,ok2 = mys2.(*MyService2)
	var myS2, ok2 = mys2.(MyService2)
	assert.IsEqual(ok2, true)
	_ = myS2.Create(dto)
}
