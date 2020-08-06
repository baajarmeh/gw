package gw

import (
	"fmt"
	"strings"
	"sync"
)

type DecoratorPoint int8

const (
	DecoratorPointActionBefore DecoratorPoint = 0
	DecoratorPointActionAfter
)

type IDecorator interface {
	Catalog() string
	Point() DecoratorPoint
	Call(ctx *Context) (friendlyMsg string, err error)
}

type DecoratorPermissionImpl struct {
	perms       []Permission
	friendlyMsg string
}

const permissionDecoratorCatalog = "permission"

var ErrPermissionDenied = fmt.Errorf("permission denied")

func (p DecoratorPermissionImpl) Catalog() string {
	return permissionDecoratorCatalog
}

func (p DecoratorPermissionImpl) Point() DecoratorPoint {
	return DecoratorPointActionBefore
}

func (p DecoratorPermissionImpl) Call(ctx *Context) (friendlyMsg string, err error) {
	s := hostServer(ctx.Context)
	if !s.permissionManager.HasPermission(ctx.User, p.perms...) {
		return p.friendlyMsg, ErrPermissionDenied
	}
	return "", nil
}

type PermissionDecoratorList struct {
	locker         sync.Mutex
	items          []IDecorator
	permDecorators map[string]IDecorator
}

func (p PermissionDecoratorList) Administration() IDecorator {
	return p.permDecorators["Administration"]
}

func (p PermissionDecoratorList) Creation() IDecorator {
	return p.permDecorators["Creation"]
}

func (p PermissionDecoratorList) Modification() IDecorator {
	return p.permDecorators["Modification"]
}

func (p PermissionDecoratorList) Deletion() IDecorator {
	return p.permDecorators["Deletion"]
}

func (p PermissionDecoratorList) ReadAll() IDecorator {
	return p.permDecorators["ReadAll"]
}

func (p PermissionDecoratorList) ReadDetail() IDecorator {
	return p.permDecorators["ReadDetail"]
}

func (p PermissionDecoratorList) Has(name string) bool {
	var _, ok = p.permDecorators[name]
	return ok
}

func (p *PermissionDecoratorList) All() []IDecorator {
	p.locker.Lock()
	defer p.locker.Unlock()
	if p.items == nil {
		var items []IDecorator
		for _, v := range p.permDecorators {
			items = append(items, v)
		}
		p.items = items
	}
	var list []IDecorator
	copy(list, p.items)
	return list
}

// NewAllPermDecorator returns a PermissionDecoratorList, that has
// Administration,Creation,Deletion,Modification,RealAll,ReadDetail Permission
func NewAllPermDecorator(resource string) PermissionDecoratorList {
	var pdList = NewCrudPermDecorator(resource)
	pdList.permDecorators["Administration"] = NewAdministrationPermDecorator(resource)
	return pdList
}

// NewResourceCreationPermDecorator returns a resources of has perm Permission Decorator.
func NewPermDecorator(perm, resource string, suffix ...string) IDecorator {
	kn := fmt.Sprintf("%s%s%sPermission", perm, resource, suffix)
	desc := fmt.Sprintf("A %s%s permssion for %s", perm, suffix, resource)
	return NewPermissionDecorator(NewPermSameKeyName(kn, desc))
}

// NewAdministrationPermDecorator returns a resources of has full(administration) Permission object.
func NewAdministrationPermDecorator(resource string) IDecorator {
	return NewPermDecorator("Administration", resource)
}

// NewResourceCreationPermDecorator returns a resources of has Creation Permission object.
func NewCreationPermDecorator(resource string) IDecorator {
	return NewPermDecorator("Creation", resource)
}

// NewModificationPermDecorator returns a resources of has Modification Permission object.
func NewModificationPermDecorator(resource string) IDecorator {
	return NewPermDecorator("Modification", resource)
}

// NewDeletionPermDecorator returns a has Deletion Permission of resources.
func NewDeletionPermDecorator(resource string) IDecorator {
	return NewPermDecorator("Deletion", resource)
}

// NewReadDetailPermDecorator returns a has Read detail Permission of resources.
func NewReadDetailPermDecorator(resource string) IDecorator {
	return NewPermDecorator("Read", resource, "Detail")
}

// NewReadAllPermDecorator returns a has Read all/pager List Permission of resources.
func NewReadAllPermDecorator(resource string) IDecorator {
	return NewPermDecorator("Read", resource, "All")
}

// NewCrudPermDecorator returns
// A resources Create, Delete,
func NewCrudPermDecorator(resource string) PermissionDecoratorList {
	var pdList = PermissionDecoratorList{
		permDecorators: make(map[string]IDecorator),
	}
	pdList.permDecorators["ReadAll"] = NewReadAllPermDecorator(resource)
	pdList.permDecorators["Creation"] = NewCreationPermDecorator(resource)
	pdList.permDecorators["Deletion"] = NewDeletionPermDecorator(resource)
	pdList.permDecorators["Modification"] = NewModificationPermDecorator(resource)
	pdList.permDecorators["ReadDetail"] = NewReadDetailPermDecorator(resource)
	return pdList
}

func NewPermissionDecorator(perms ...Permission) IDecorator {
	names := make([]string, len(perms))
	for idx := 0; idx < len(perms); idx++ {
		names[idx] = perms[idx].Name
	}
	friendlyMsg := fmt.Sprintf("Permission Deined, needs:(%s)", strings.Join(names, "|"))
	return &DecoratorPermissionImpl{
		perms:       perms,
		friendlyMsg: friendlyMsg,
	}
}

// helpers
func filterDecorator(filter func(d IDecorator) bool, decorators ...IDecorator) []IDecorator {
	var result []IDecorator
	for _, dc := range decorators {
		if filter(dc) {
			result = append(result, dc)
		}
	}
	return result
}
