package gw

import (
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type Dto1 struct {
}

type MyUser struct {
}

type IMyService interface {
	Create(dto Dto1) error
}

type MyService struct {
	User  *gorm.DB
	typer IMyService
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
func (MyService2) New(store IStore, service IMyService) *MyService2 {
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
func (MyService3) New(store IStore, service IMyService, service2 MyService2) *MyService3 {
	var myService = MyService3{}
	myService.User = store.GetDbStore().Model(MyUser{})
	myService.MyService1 = service
	myService.MyService2 = service2
	return &myService
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
func (MyService4) New(store IStore, service IMyService, service2 MyService2, service3 *MyService3) MyService4 {
	var myService = MyService4{}
	myService.User = store.GetDbStore().Model(MyUser{})
	myService.MyService1 = service
	myService.MyService2 = service2
	myService.MyService3 = service3
	return myService
}

type MyServices struct {
	MyService1 IMyService
	MyService2 MyService2
	MyService3 *MyService3
	MyService4 MyService4
}

func (s MyServices) New(service IMyService, service2 MyService2,
	service3 *MyService3, service4 MyService4) MyServices {
	s.MyService1 = service
	s.MyService2 = service2
	s.MyService3 = service3
	s.MyService4 = service4
	return s
}

var diTester IDIProvider
var diTesterStore IStore
var diTesterBuiltinComponent BuiltinComponent
var myServicesTyper reflect.Type

func init() {
	var server = NewTesterServer("default-tester")
	diTester = server.DIProvider
	var myS1Impl MyService
	var myS2Impl MyService2
	var myS3Impl MyService3
	var myS4Impl MyService4
	var myServices MyServices
	myServicesTyper = reflect.TypeOf(myServices)
	diTester.Register(myS1Impl, myS2Impl, myS3Impl, myS4Impl, myServices)
	e, d := diTester.ResolveByTyper(BuiltinComponentTyper)
	if e != nil {
		panic(e)
	}
	diTesterBuiltinComponent = d.(BuiltinComponent)
}

func BenchmarkDefaultDIProviderImpl_ResolveByTyperWithState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = diTester.ResolveByTyperWithState(diTesterStore, myServicesTyper)
	}
}

func TestDefaultDIProviderImpl_BuiltinComponent(t *testing.T) {
	assert2.True(t, diTesterBuiltinComponent.Store != nil)
	assert2.True(t, diTesterBuiltinComponent.IDGenerator != nil)
}

// FIXME(OceanHo): go test ./*.go can not pass ?
func TestDefaultDIProviderImpl(t *testing.T) {
	var myServices MyServices
	err := diTester.ResolveByObjectTyper(&myServices)
	assert.IsEqual(err, nil)
	assert.NotEqual(t, myServices, nil)
	var dto Dto1
	_ = myServices.MyService1.Create(dto)
	_ = myServices.MyService2.Create(dto)
	_ = myServices.MyService3.Create(dto)
	_ = myServices.MyService4.Create(dto)
}
