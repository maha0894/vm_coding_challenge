package db

import (
	"fmt"
	"log"
	"time"
	"vm_coding_challenge/config"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

const MaxDatabaseConnectionAttempts int = 10

func Setup() (err error) {
	migrateConf := &goose.DBConf{
		MigrationsDir: "DB/migrations",
		Env:           "production",
		Driver: &goose.DBDriver{
			Name:    "mysql",
			OpenStr: config.Conf.DBPath,
			Import:  "github.com/go-sql-driver/mysql",
			Dialect: &goose.MySqlDialect{},
		},
	}
	// Get the latest possible migration
	latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)
	if err != nil {
		fmt.Println("DB error: ", err)
		return err
	}

	// Open our database connection
	i := 0
	for {
		DB, err = gorm.Open("mysql", config.Conf.DBPath)
		if err == nil {
			break
		}
		if err != nil && i >= MaxDatabaseConnectionAttempts {
			fmt.Println("DB error: ", err)
			return err
		}
		i++
		fmt.Println("waiting for database to be up...")
		time.Sleep(5 * time.Second)
	}
	DB.LogMode(false)
	DB.SetLogger(log.Logger)
	DB.DB().SetMaxOpenConns(1)
	if err != nil {
		fmt.Println("DB error: ", err)
		return err
	}
	// Migrate up to the latest version
	err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, DB.DB())
	if err != nil {
		fmt.Println("DB error: ", err)
		return err
	}
	return
}
