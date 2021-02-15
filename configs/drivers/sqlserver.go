package drivers

import (
	"fmt"
	"log"

	driver "gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqlServer struct {
}

func (d *SqlServer) Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB {
	var db *gorm.DB
	var err error

	conn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, host, port, dbname)
	if debug {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{
			SkipDefaultTransaction: true,
		})
	} else {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Silent),
			SkipDefaultTransaction: true,
		})
	}

	if err != nil {
		log.Printf("Gorm SQL Server: %+v \n", err)
		panic(err)
	}

	return db
}
