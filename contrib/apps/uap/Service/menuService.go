package Service

import (
	"fmt"
	"github.com/oceanho/gw"
	"github.com/oceanho/gw/contrib/apps/uap/Dto"
)

type IMenuService interface {
	Create(dto *Dto.BatchCreateMenuDto) error
	Modify(dto *Dto.Menu) error
	Delete(id uint64) error
	Query(id uint64) (error, *Dto.Menu)
	QueryByKey(key string) (error, *Dto.Menu)
	QueryList(expr gw.QueryExpr) (error, *MenuPagerResult)
	QueryListByApp(appId uint64, expr gw.QueryExpr) (error, *MenuPagerResult)
}

type MenuPagerResult struct {
	gw.PagerResult
	Data []Dto.Menu `json:"data"`
}

type MenuService struct {
	gw.BuiltinComponent
}

// DI
func (ms MenuService) New(bc gw.BuiltinComponent) IMenuService {
	ms.BuiltinComponent = bc
	return &ms
}

func (ms MenuService) Create(dto *Dto.BatchCreateMenuDto) error {
	db := ms.Store.GetDbStore()
	var appInfo = ms.AppManager.QueryByName(dto.App)
	if appInfo == nil {
		return fmt.Errorf("system not found app, key=%s", dto.App)
	}
	return db.Error
}

func (ms MenuService) Modify(dto *Dto.Menu) error {
	panic("implement me")
}

func (ms MenuService) Delete(id uint64) error {
	panic("implement me")
}

func (ms MenuService) Query(id uint64) (error, *Dto.Menu) {
	panic("implement me")
}

func (ms MenuService) QueryByKey(key string) (error, *Dto.Menu) {
	panic("implement me")
}

func (ms MenuService) QueryList(expr gw.QueryExpr) (error, *MenuPagerResult) {
	panic("implement me")
}

func (ms MenuService) QueryListByApp(appId uint64, expr gw.QueryExpr) (error, *MenuPagerResult) {
	panic("implement me")
}
