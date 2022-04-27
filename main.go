package main

import (
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
	"xinxi/pkg"
)

type Config struct {
	Dsn string `yaml:"dsn"`
	Redis string `yaml:"redis"`
	PinTuan *pkg.Cfg `yaml:"pin_tuan"`
}

func main() {
	n := time.Now()
	nDay := n.Format("2006-01-02 15:04")
	log.Println("运行时间:" + nDay)
	b, err := os.ReadFile("./conf.yaml")
	if err != nil {
		log.Println(err)
	}
	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Println(err)
	}

	db, err := gorm.Open(mysql.Open(c.Dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(20)
	sqlDb.SetMaxOpenConns(30)

	p := pkg.New(c.PinTuan, db)
	p.TimeTicker()
	p.Run()
}
