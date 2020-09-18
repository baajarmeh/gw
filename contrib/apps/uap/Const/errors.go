package Const

import "fmt"

var (
	ErrorNonUserCannotCreationResource      = fmt.Errorf("non-user cannot be to allow create this resource")
	ErrorTenancyCannotCreationAdminResource = fmt.Errorf("this operation are denied")
)
