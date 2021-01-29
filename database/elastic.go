package database

import (
	"github.com/organization-service/goorg/internal"
	apmmysql "go.elastic.co/apm/module/apmgormv2/driver/mysql"
	apmpostgres "go.elastic.co/apm/module/apmgormv2/driver/postgres"
	apmsqlite "go.elastic.co/apm/module/apmgormv2/driver/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func dialect(environment *environment) gorm.Dialector {
	switch environment.Dialect {
	case "postgres":
		if internal.ApmName == internal.Elastic {
			return apmpostgres.Open(environment.DataSource)
		}
		return postgres.Open(environment.DataSource)
	case "sqlite3":
		if internal.ApmName == internal.Elastic {
			return apmsqlite.Open(environment.DataSource)
		}
		return sqlite.Open(environment.DataSource)
	case "mysql":
		if internal.ApmName == internal.Elastic {
			return apmmysql.Open(environment.DataSource)
		}
		return mysql.Open(environment.DataSource)
	default:
		return nil
	}
}
