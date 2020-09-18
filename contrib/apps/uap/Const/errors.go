package Const

import "fmt"

var (
	ErrorNonUserCannotCreationResource      = fmt.Sprintf("Non-user cannot be to allow create this resource.")
	ErrorTenancyCannotCreationAdminResource = fmt.Sprintf("This operation are denied.")
)
