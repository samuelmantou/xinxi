package pkg

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
	"xinxi/pkg/model"
)

func TestA(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/xinxi?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	var c model.ConfigData
	if err := db.Where("`root_key` = 'xinxi'  AND `key` = ?", "xinxi2_pin_tuan_start").Find(&c).Error; err != nil {
		log.Println(err)
	}
	log.Println(c)
}
