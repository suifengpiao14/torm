package sourceprovider

import (
	"github.com/suifengpiao14/sqlexec"
	"github.com/suifengpiao14/sshmysql"
)

var DriverName = "mysql"

func NewDBProvider(dbConfig sqlexec.DBConfig, sshConfig *sshmysql.SSHConfig) (dbProvider *sqlexec.ExecutorSQL) {
	dbProvider = sqlexec.NewExecutorSQL(dbConfig, sshConfig)

	return dbProvider
}
