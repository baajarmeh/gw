package gw

type IUserManager interface {
	Create(user *User) error
	Modify(user User) error
	Delete(tenantId, userId uint64) error
	Query(tenantId, userId uint64) (User, error)
	QueryList(tenantId uint64, expr PagerExpr, total int64, out []User) error
}

func DefaultUserManager(state ServerState) IUserManager {
	return DefaultUserManagerImpl{
		state: state,
	}
}

// DefaultUserManagerImpl ...
type DefaultUserManagerImpl struct {
	state ServerState
}

func (d DefaultUserManagerImpl) Create(user *User) error {
	pm := d.state.PermissionManager()
	return pm.GrantToUser(user.Id, user.Permissions...)
}

func (d DefaultUserManagerImpl) Modify(user User) error {
	panic("implement me")
}

func (d DefaultUserManagerImpl) Delete(tenantId, userId uint64) error {
	panic("implement me")
}

func (d DefaultUserManagerImpl) Query(tenantId, userId uint64) (User, error) {
	panic("implement me")
}

func (d DefaultUserManagerImpl) QueryList(tenantId uint64, expr PagerExpr, total int64, out []User) error {
	panic("implement me")
}
