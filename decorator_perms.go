package gw

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const permissionDecoratorCatalog = "gw_framework_permission"

var (
	ErrUnauthorized        = fmt.Errorf("has no credentitals")
	ErrInternalServerError = fmt.Errorf("server internal error")
	ErrBadRequest          = fmt.Errorf("bad request")
	ErrNotFoundRequest     = fmt.Errorf("not found")
	ErrPermissionDenied    = fmt.Errorf("permission denied")
)

type PermissionDecorator struct {
	resName        string
	locker         sync.Mutex
	permDecorators map[string]Decorator
}

func (p *PermissionDecorator) Administration() Decorator {
	return p.permDecorators["Administration"]
}

func (p *PermissionDecorator) Creation() Decorator {
	return p.permDecorators["Creation"]
}

func (p *PermissionDecorator) Modification() Decorator {
	return p.permDecorators["Modification"]
}

func (p *PermissionDecorator) Deletion() Decorator {
	return p.permDecorators["Deletion"]
}

func (p *PermissionDecorator) ReadAll() Decorator {
	return p.permDecorators["ReadAll"]
}

func (p *PermissionDecorator) ReadDetail() Decorator {
	return p.permDecorators["ReadDetail"]
}

func (p *PermissionDecorator) Has(name string) bool {
	var _, ok = p.permDecorators[name]
	return ok
}

func (p *PermissionDecorator) Extend(op string) *PermissionDecorator {
	if p.Has(op) {
		return p
	}
	p.permDecorators[op] = NewPermDecorator(op, p.resName)
	return p
}

func (p *PermissionDecorator) Merge(decorators ...*PermissionDecorator) {
	p.locker.Lock()
	defer p.locker.Unlock()
	for _, d := range decorators {
		for k, v := range d.permDecorators {
			p.permDecorators[k] = v
		}
	}
}

func (p *PermissionDecorator) Permissions() []*Permission {
	var perms []*Permission
	for _, item := range p.permDecorators {
		item := item
		if perm, ok := item.MetaData.([]*Permission); ok {
			perms = append(perms, perm...)
		}
	}
	return perms
}

func (p *PermissionDecorator) All() []Decorator {
	p.locker.Lock()
	defer p.locker.Unlock()
	var items []Decorator
	for _, v := range p.permDecorators {
		items = append(items, v)
	}
	return items
}

// NewPermAllDecorator returns a PermissionDecoratorList, that has
// Administration,Creation,Deletion,Modification,RealAll,ReadDetail Permission
func NewPermAllDecorator(resName string) *PermissionDecorator {
	var pdList = NewCrudPermDecorator(resName)
	pdList.permDecorators["All"] = NewAdministrationPermDecorator(resName)
	return pdList
}

// NewPermDecorator returns a resNames of has perm Permission Decorator.
func NewPermDecorator(op, resName string) Decorator {
	kn := fmt.Sprintf("%s.%s", resName, op)
	desc := fmt.Sprintf("Define a %s permssion on resource %s.", op, resName)
	return NewPermissionDecorator(NewPermSameKeyName(kn, desc))
}

// NewAdministrationPermDecorator returns a resNames of has full(administration) Permission object.
func NewAdministrationPermDecorator(resName string) Decorator {
	return NewPermDecorator("All", resName)
}

// NewCreationPermDecorator returns a resNames of has Creation Permission object.
func NewCreationPermDecorator(resName string) Decorator {
	return NewPermDecorator("Create", resName)
}

// NewModificationPermDecorator returns a resNames of has Modification Permission object.
func NewModificationPermDecorator(resName string) Decorator {
	return NewPermDecorator("Modify", resName)
}

// NewDeletionPermDecorator returns a has Deletion Permission of resNames.
func NewDeletionPermDecorator(resName string) Decorator {
	return NewPermDecorator("Delete", resName)
}

// NewReadPermDecorator returns a has Read detail Permission of resNames.
func NewReadPermDecorator(resName string) Decorator {
	return NewPermDecorator("Read", resName)
}

// NewQueryPermDecorator returns a has Read all/pager List Permission of resNames.
func NewQueryPermDecorator(resName string) Decorator {
	return NewPermDecorator("Query", resName)

}

// NewCrudPermDecorator returns
// A resNames Create, Delete,
func NewCrudPermDecorator(resName string) *PermissionDecorator {
	var pdList = &PermissionDecorator{
		resName:        resName,
		permDecorators: make(map[string]Decorator),
	}
	pdList.permDecorators["Query"] = NewQueryPermDecorator(resName)
	pdList.permDecorators["Create"] = NewCreationPermDecorator(resName)
	pdList.permDecorators["Delete"] = NewDeletionPermDecorator(resName)
	pdList.permDecorators["Modify"] = NewModificationPermDecorator(resName)
	pdList.permDecorators["Read"] = NewReadPermDecorator(resName)
	return pdList
}

func NewBeforeCustomFuncDecorator(handler DecoratorHandler) Decorator {
	return Decorator{
		MetaData: nil,
		Before:   handler,
		After:    nil,
	}
}

func NewPermissionDecorator(perms ...*Permission) Decorator {
	names := make([]string, len(perms))
	for idx := 0; idx < len(perms); idx++ {
		names[idx] = perms[idx].Name
	}
	msg := fmt.Sprintf("Permission Deined, need:(%s)", strings.Join(names, "|"))
	return Decorator{
		MetaData: perms,
		Before: func(c *Context) (status int, err error, payload interface{}) {
			s := c.Server()
			if s.PermissionManager.Checker().Check(c.User(), perms...) {
				return 0, nil, nil
			}
			return http.StatusForbidden, ErrPermissionDenied, msg
		},
		After: nil,
	}
}
