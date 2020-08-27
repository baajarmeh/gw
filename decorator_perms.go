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

func (p *PermissionDecorator) Merge(decorators ...*PermissionDecorator) {
	p.locker.Lock()
	defer p.locker.Unlock()
	for _, d := range decorators {
		for k, v := range d.permDecorators {
			p.permDecorators[k] = v
		}
	}
}

func (p *PermissionDecorator) Permissions() []Permission {
	var perms []Permission
	for _, item := range p.permDecorators {
		item := item
		if perm, ok := item.MetaData.([]Permission); ok {
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
func NewPermAllDecorator(resource string) *PermissionDecorator {
	var pdList = NewCrudPermDecorator(resource)
	pdList.permDecorators["Administration"] = NewAdministrationPermDecorator(resource)
	return pdList
}

// NewPermDecorator returns a resources of has perm Permission Decorator.
func NewPermDecorator(perm, resource string) Decorator {
	return NewPermDecoratorWithSuffix(perm, resource, "")
}

// NewPermDecoratorWithSuffix returns a resources of has perm Permission Decorator.
func NewPermDecoratorWithSuffix(perm, resource string, suffix string) Decorator {
	kn := fmt.Sprintf("%s%s%sPermission", perm, resource, suffix)
	desc := fmt.Sprintf("A %s%s permssion for %s", perm, suffix, resource)
	return NewPermissionDecorator(NewPermSameKeyName(kn, desc))
}

// NewPermDecoratorWithPrefix returns a resources of has perm Permission Decorator.
func NewPermDecoratorWithPrefix(perm, resource string, prefix string) Decorator {
	kn := fmt.Sprintf("%s%s%sPermission", perm, prefix, resource)
	desc := fmt.Sprintf("A %s%s permssion for %s", perm, prefix, resource)
	return NewPermissionDecorator(NewPermSameKeyName(kn, desc))
}

// NewAdministrationPermDecorator returns a resources of has full(administration) Permission object.
func NewAdministrationPermDecorator(resource string) Decorator {
	return NewPermDecorator("Administration", resource)
}

// NewCreationPermDecorator returns a resources of has Creation Permission object.
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
	return NewPermDecoratorWithSuffix("Read", resource, "Detail")
}

// NewReadAllPermDecorator returns a has Read all/pager List Permission of resources.
func NewReadAllPermDecorator(resource string) Decorator {
	return NewPermDecoratorWithPrefix("Read", resource, "All")
}

// NewCrudPermDecorator returns
// A resources Create, Delete,
func NewCrudPermDecorator(resource string) *PermissionDecorator {
	var pdList = &PermissionDecorator{
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
	msg := fmt.Sprintf("Permission Deined, need:(%s)", strings.Join(names, "|"))
	return Decorator{
		MetaData: perms,
		Before: func(c *Context) (status int, err error, payload interface{}) {
			s := c.HostServer()
			if s.PermissionManager.Checker().Check(c.User(), perms...) {
				return 0, nil, nil
			}
			return http.StatusForbidden, ErrPermissionDenied, msg
		},
		After: nil,
	}
}
