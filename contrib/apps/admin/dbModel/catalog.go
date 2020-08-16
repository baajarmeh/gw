package dbModel

type Catalog struct {
}

func (Catalog) tableName() string {
	return "catalog"
}

func (t Catalog) TableName() string {
	return GetGwContribAdminTableName(t.tableName())
}
