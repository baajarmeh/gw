package gw

import "fmt"

var (
	errNoAuthorization  = fmt.Errorf("no authorization")
	errPassportOrSecret = fmt.Errorf("passport or secret error")
)
