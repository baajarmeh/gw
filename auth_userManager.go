package gw

type defaultUser struct {
	User
	secret string
}

// DefaultUserManagerImpl ...
type DefaultUserManagerImpl struct {
	state *ServerState
}

func (d DefaultUserManagerImpl) Create(user *User) error {
	return nil
}

func (d DefaultUserManagerImpl) Modify(user User) error {
	return nil
}

func (d DefaultUserManagerImpl) Delete(tenantId, userId uint64) error {
	return nil
}

func (d DefaultUserManagerImpl) Query(tenantId, userId uint64) (User, error) {
	return EmptyUser, nil
}

func (d DefaultUserManagerImpl) QueryByUser(tenantId uint64, user, password string) (User, error) {
	return EmptyUser, nil
}

func (d DefaultUserManagerImpl) QueryByAKS(tenantId uint64, accessKey, accessSecret string) (User, error) {
	return EmptyUser, nil
}

func (d DefaultUserManagerImpl) QueryList(tenantId uint64, expr PagerExpr, total int64, out []User) error {
	return nil
}
