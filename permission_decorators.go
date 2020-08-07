package gw

import (
	"fmt"
	"strings"
	"sync"
)

type DecoratorPermissionImpl struct {
	perms       []Permission
	friendlyMsg string
}

const permissionDecoratorCatalog = "permission"

var ErrPermissionDenied = fmt.Errorf("permission denied")

func (p DecoratorPermissionImpl) Catalog() string {
	return permissionDecoratorCatalog
}

func (p DecoratorPermissionImpl) OnBeforeCall(ctx *Context) (friendlyMsg string, err error) {
	s := GetHostServer(ctx.Context)
	if !s.PermissionManager.HasPermission(ctx.User, p.perms...) {
		return p.friendlyMsg, ErrPermissionDenied
	}
	return "", nil
}

func (p DecoratorPermissionImpl) OnAfterCall(ctx *Context) (friendlyMsg string, err error) {
	return "", nil
}

type PermissionDecoratorList struct {
	locker         sync.Mutex
	items          []Decorator
	permDecorators map[string]Decorator
}

func (p PermissionDecoratorList) Administration() Decorator {
	return p.permDecorators["Administration"]
}

func (p PermissionDecoratorList) Creation() Decorator {
	return p.permDecorators["Creation"]
}

func (p PermissionDecoratorList) Modification() Decorator {
	return p.permDecorators["Modification"]
}

func (p PermissionDecoratorList) Deletion() Decorator {
	return p.permDecorators["Deletion"]
}

func (p PermissionDecoratorList) ReadAll() Decorator {
	return p.permDecorators["ReadAll"]
}

func (p PermissionDecoratorList) ReadDetail() Decorator {
	return p.permDecorators["ReadDetail"]
}

func (p PermissionDecoratorList) Has(name string) bool {
	var _, ok = p.permDecorators[name]
	return ok
}

func (p *PermissionDecoratorList) All() []Decorator {
	p.locker.Lock()
	defer p.locker.Unlock()
	if p.items == nil {
		var items []Decorator
		for _, v := range p.permDecorators {
			items = append(items, v)
		}
		p.items = items
	}
	var list = make([]Decorator, len(p.items))
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
func NewPermDecorator(perm, resource string, suffix ...string) Decorator {
	kn := fmt.Sprintf("%s%s%sPermission", perm, resource, suffix)
	desc := fmt.Sprintf("A %s%s permssion for %s", perm, suffix, resource)
	return NewPermissionDecorator(NewPermSameKeyName(kn, desc))
}

// NewAdministrationPermDecorator returns a resources of has full(administration) Permission object.
func NewAdministrationPermDecorator(resource string) Decorator {
	return NewPermDecorator("Administration", resource)
}

// NewResourceCreationPermDecorator returns a resources of has Creation Permission object.
func NewCreationPermDecorator(resource string) Decorator {
	return NewPermDecorator("Creation", resource)
}

// NewModificationPermDecorator returns a resources of has Modification Permission object.
func NewModificationPermDecorator(resource string) Decorator {
	return NewPermDecorator("Modification", resource)
}

// NewDeletionPermDecorator returns a has Deletion Permission of resources.
func NewDeletionPermDecorator(resource string) Decorator {
	return NewPermDecorator("Deletion", resource)
}

// NewReadDetailPermDecorator returns a has Read detail Permission of resources.
func NewReadDetailPermDecorator(resource string) Decorator {
	return NewPermDecorator("Read", resource, "Detail")
}

// NewReadAllPermDecorator returns a has Read all/pager List Permission of resources.
func NewReadAllPermDecorator(resource string) Decorator {
	return NewPermDecorator("Read", resource, "All")
}

// NewCrudPermDecorator returns
// A resources Create, Delete,
func NewCrudPermDecorator(resource string) PermissionDecoratorList {
	var pdList = PermissionDecoratorList{
		permDecorators: make(map[string]Decorator),
	}
	pdList.permDecorators["ReadAll"] = NewReadAllPermDecorator(resource)
	pdList.permDecorators["Creation"] = NewCreationPermDecorator(resource)
	pdList.permDecorators["Deletion"] = NewDeletionPermDecorator(resource)
	pdList.permDecorators["Modification"] = NewModificationPermDecorator(resource)
	pdList.permDecorators["ReadDetail"] = NewReadDetailPermDecorator(resource)
	return pdList
}

func NewPermissionDecorator(perms ...Permission) Decorator {
	names := make([]string, len(perms))
	for idx := 0; idx < len(perms); idx++ {
		names[idx] = perms[idx].Name
	}
	msg := fmt.Sprintf("Permission Deined, needs:(%s)", strings.Join(names, "|"))
	return Decorator{
		MetaData: perms,
		Before: func(c *Context) (friendlyMsg string, err error) {
			s := GetHostServer(c.Context)
			if s.PermissionManager.HasPermission(c.User, perms...) {
				return "", nil
			}
			return msg, ErrPermissionDenied
		},
		After: nil,
	}
}
