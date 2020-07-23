package gw

import "fmt"

var (
	NoAuthorizationError error = fmt.Errorf("no authorization")
	PassportOrSecretError error = fmt.Errorf("passport or secret error")
)