package Db

import "fmt"

const tablePrefix = "gw_uap"

func getTableName(name string) string {
	return fmt.Sprintf("%s_%s", tablePrefix, name)
}
