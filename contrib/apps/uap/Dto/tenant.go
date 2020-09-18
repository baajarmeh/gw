package Dto

//go:generate gomodifytags -w -file ./tenant.go -add-tags json -transform snakecase -all

type TenantDto struct {
	Name string `json:"name"`
}
