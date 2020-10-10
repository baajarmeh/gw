package Const

import "fmt"

var (
	ErrorNonUserCannotModifyResource        = fmt.Errorf("non-user cannot be to allow modify this resource")
	ErrorNonUserCannotCreationResource      = fmt.Errorf("non-user cannot be to allow create this resource")
	ErrorTenancyCannotCreationAdminResource = fmt.Errorf("permission denied of current operation")
)
