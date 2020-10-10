package Const

import "fmt"

var (
	ErrorNonUserCannotModifyResource   = fmt.Errorf("non-user cannot be to allow modify this resource")
	ErrorNonUserCannotCreationResource = fmt.Errorf("non-user cannot be to allow create this resource")
	ErrorNonUserHasNoOperatePermission = fmt.Errorf("non-user has no permission to operate this resource")
)
