package gw

import "github.com/go-playground/validator/v10"

type Validator struct {
	Name    string
	Message string
	Handler validator.Func
}

func NewValidator(name, message string, handler validator.Func) *Validator {
	return &Validator{
		Name:    name,
		Message: message,
		Handler: handler,
	}
}