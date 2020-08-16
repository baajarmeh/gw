package dbModel

import "fmt"

const (
	GwContribTableNamePrefix      = "gw_contrib_"
	GwContribAdminTableNamePrefix = "gw_contrib_admin_"
)

func GetGwContribTableName(tableName string) string {
	return fmt.Sprintf("%s%s", GwContribTableNamePrefix, tableName)
}

func GetGwContribAdminTableName(tableName string) string {
	return fmt.Sprintf("%s%s", GwContribAdminTableNamePrefix, tableName)
}
