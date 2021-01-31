package database

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/organization-service/goorg/internal"
	apmmysql "go.elastic.co/apm/module/apmgormv2/driver/mysql"
	apmpostgres "go.elastic.co/apm/module/apmgormv2/driver/postgres"
	apmsqlite "go.elastic.co/apm/module/apmgormv2/driver/sqlite"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	instanceDB struct {
		connectionReadWrite *gorm.DB
		connectionReadOnly  *gorm.DB
		logMode             bool
	}
	environment struct {
		Dialect    string `yaml:"dialect"`
		DataSource string `yaml:"datasource"`
	}
)

func getEnv(name, _default string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return _default
}

func (db *instanceDB) ReadWriteConnection(c context.Context) *gorm.DB {
	_db := db.connectionReadWrite.WithContext(c)
	if db.logMode {
		_db = _db.Debug()
	}
	return _db
}

func (db *instanceDB) ReadWriteConnectionObject(c context.Context) *sql.DB {
	sql, _ := db.ReadWriteConnection(c).DB()
	return sql
}

func (db *instanceDB) ReadOnlyConnection(c context.Context) *gorm.DB {
	_db := db.connectionReadOnly.WithContext(c)
	if db.logMode {
		_db = _db.Debug()
	}
	return _db
}

func (db *instanceDB) ReadOnlyConnectionObject(c context.Context) *sql.DB {
	sql, _ := db.ReadOnlyConnection(c).DB()
	return sql
}

func (db *instanceDB) LogMode(set bool) {
	db.logMode = set
}

func New() IDriver {
	apmName := getEnv("APM_NAME", "")
	if db != nil {
		return db
	}
	masterEnv, err := readConfig(Master)
	if err != nil {
		panic(err)
	}
	var slaveEnv *environment
	if slaveEnv, err = readConfig(Slave); err != nil {
		if errNotFoundKey == err {
			master := *masterEnv
			slaveEnv = &master
		} else {
			panic(err)
		}
	}

	fnDialect := func(environment *environment) gorm.Dialector {
		switch environment.Dialect {
		case "postgres":
			if apmName == internal.Elastic {
				return apmpostgres.Open(environment.DataSource)
			}
			return postgres.Open(environment.DataSource)
		case "sqlite3":
			if apmName == internal.Elastic {
				return apmsqlite.Open(environment.DataSource)
			}
			return sqlite.Open(environment.DataSource)
		case "mysql":
			if apmName == internal.Elastic {
				return apmmysql.Open(environment.DataSource)
			}
			return mysql.Open(environment.DataSource)
		}
		return nil
	}
	masterDB, err := gorm.Open(fnDialect(masterEnv), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	slaveDB, err := gorm.Open(fnDialect(slaveEnv), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	setConnectionPool := func(db *gorm.DB) {
		sqldb, _ := db.DB()
		sqldb.SetMaxIdleConns(10)
		sqldb.SetMaxOpenConns(100)
		sqldb.SetConnMaxLifetime(time.Hour)
	}

	// ログ出力
	masterDB.Logger.LogMode(logger.Info)
	slaveDB.Logger.LogMode(logger.Info)
	if apmName == internal.Elastic {
		addGormCallbacks(masterDB)
		addGormCallbacks(slaveDB)
	}
	// コネクションプーリング設定
	setConnectionPool(masterDB)
	setConnectionPool(slaveDB)
	db = &instanceDB{
		connectionReadWrite: masterDB,
		connectionReadOnly:  slaveDB,
	}
	return db
}

func readConfigs() (map[string]*environment, error) {
	if configs != nil {
		return configs, nil
	}
	configDir := getEnv("DB_CONFIG_DIR", "")
	configFile := path.Join(configDir, "dbconfig.yml")
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	configs = make(map[string]*environment)
	err = yaml.Unmarshal(file, configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func readConfig(key string) (*environment, error) {
	config, err := readConfigs()
	if err != nil {
		return nil, err
	}
	if conf, ok := config[key]; ok {
		conf.DataSource = os.ExpandEnv(conf.DataSource)
		return conf, nil
	} else {
		return nil, errNotFoundKey
	}
}
