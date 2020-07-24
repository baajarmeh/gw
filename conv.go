package gw

type IToDto interface {
	ToDto(dto interface{}) error
}

type IToModel interface {
	ToModel(model interface{}) error
}
