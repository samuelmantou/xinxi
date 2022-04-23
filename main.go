package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"runtime"
	"xinxi/pkg"
)

func main() {
	var dsn string
	if runtime.GOOS == "darwin" {
		dsn = "root:@tcp(127.0.0.1:3306)/xinxi?charset=utf8mb4&parseTime=True&loc=Local"
	}else if runtime.GOOS == "windows" {
		dsn = "root:@tcp(127.0.0.1:3306)/x?charset=utf8mb4&parseTime=True&loc=Local"
	}else{
		dsn = ""
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(20)
	sqlDb.SetMaxOpenConns(30)

	cfg := &pkg.Cfg{
		Change: 5,
		Round: 30,
		Winner: 5,
		Insert: 2,
		Start: "07:00",
		End: "20:00",
	}
	p := pkg.New(cfg, db)
	p.TimeTicker()
	p.Run()
}
