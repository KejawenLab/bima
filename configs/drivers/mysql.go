package drivers

import (
	"fmt"
	"log"
	"os"
	"time"

	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
}

func (d *Mysql) Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB {
	var db *gorm.DB
	var err error

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", user, password, host, port, dbname)
	if debug {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: time.Second,
					LogLevel:      logger.Info,
					Colorful:      false,
				},
			),
		})
	} else {
		db, err = gorm.Open(driver.Open(conn), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: 200 * time.Microsecond,
					LogLevel:      logger.Warn,
					Colorful:      false,
				},
			),
		})
	}

	if err != nil {
		log.Printf("Gorm MySQL: %+v \n", err)
		panic(err)
	}

	return db
}
