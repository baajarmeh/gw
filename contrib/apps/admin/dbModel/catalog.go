package dbModel

import "fmt"

type Catalog struct {
}

func (Catalog) tableName() string {
	return "catalog"
}

func (t Catalog) TableName() string {
	return fmt.Sprintf("%s%s", tableNamePrefix, t.tableName())
}
