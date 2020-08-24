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
	IMyTestServiceName string = "github.com/oceanho/gw.IMyService"
	MyTestService2Name string = "github.com/oceanho/gw.MyService2"
	MyTestService3Name string = "github.com/oceanho/gw.MyService3"
	MyTestService4Name string = "github.com/oceanho/gw.MyService4"
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
func (MyService2) New(serverState ServerState, store IStore, service IMyService) *MyService2 {
	var myService = MyService2{}
	myService.User = store.GetDbStore().Model(MyUser{})
	myService.MyService1 = service
	return &myService
}

func (MyService2) Destroy() {

}

type MyService3 struct {
	User       *gorm.DB
	MyService1 IMyService
	MyService2 MyService2
}

func (myService MyService3) Create(dto Dto1) error {
	return nil
}

//
// Gw DI framework dependOns
func (MyService3) New(serverState ServerState, store IStore, service IMyService, service2 MyService2) MyService3 {
	var myService = MyService3{}
	myService.User = store.GetDbStore().Model(MyUser{})
	myService.MyService1 = service
	myService.MyService2 = service2
	return myService
}

type MyService4 struct {
	User       *gorm.DB
	MyService1 IMyService
	MyService2 MyService2
	MyService3 *MyService3
}

func (myService MyService4) Create(dto Dto1) error {
	return nil
}

//
// Gw DI framework dependOns
func (MyService4) New(serverState ServerState, store IStore, service IMyService, service2 MyService2, service3 *MyService3) MyService4 {
	var myService = MyService4{}
	myService.User = store.GetDbStore().Model(MyUser{})
	myService.MyService1 = service
	myService.MyService2 = service2
	myService.MyService3 = service3
	return myService
}

func TestDefaultDIProviderImpl(t *testing.T) {
	var server = NewTesterServer()
	var state = NewServerState(server)
	var di = DefaultDIProvider(state)
	var myS1Impl MyService
	var myS2Impl MyService2
	var myS3Impl MyService3
	var myS4Impl MyService4
	di.Register(myS1Impl, myS2Impl, myS3Impl, myS4Impl)
	var myS1, ok = di.Resolve(IMyTestServiceName).(IMyService)
	assert.IsEqual(ok, true)
	var dto Dto1
	_ = myS1.Create(dto)

	// Resolve can not return a pointer Object.
	var mys2 = di.Resolve(MyTestService2Name)
	var myS2, ok2 = mys2.(*MyService2)
	assert.IsEqual(ok2, true)
	_ = myS2.Create(dto)

	var mys3 = di.Resolve(MyTestService3Name)
	var myS3, ok3 = mys3.(MyService3)
	assert.IsEqual(ok3, true)
	_ = myS3.Create(dto)

	var mys4 = di.Resolve(MyTestService4Name)
	var myS4, ok4 = mys4.(MyService4)
	assert.IsEqual(ok4, true)
	_ = myS4.Create(dto)
}
