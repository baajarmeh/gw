package Dto

import "github.com/oceanho/gw/backend/gwdb"

type Menu struct {
	ID           uint64
	Name         string
	Icon         string
	Link         string
	OpenBehavior string
	Permission   string
	Children     []Menu
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
}

type BatchCreateMenuDto struct {
	App   string
	Menus []Menu
}
