package drivers

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
}

func (_ Mysql) Connect(host string, port int, user string, password string, dbname string, debug bool) *gorm.DB {
	var db *gorm.DB
	var err error
	var dsn bytes.Buffer

	dsn.WriteString(user)
	dsn.WriteString(":")
	dsn.WriteString(password)
	dsn.WriteString("@tcp(")
	dsn.WriteString(host)
	dsn.WriteString(":")
	dsn.WriteString(strconv.Itoa(port))
	dsn.WriteString(")/")
	dsn.WriteString(dbname)
	dsn.WriteString("?charset=utf8&parseTime=true&loc=UTC")

	if debug {
		db, err = gorm.Open(mysql.Open(dsn.String()), &gorm.Config{
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
		db, err = gorm.Open(mysql.Open(dsn.String()), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: 200 * time.Millisecond,
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
