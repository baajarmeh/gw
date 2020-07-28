package dto

import "database/sql"

type  Table struct{
	Name string
	Columns []sql.ColumnType
}
