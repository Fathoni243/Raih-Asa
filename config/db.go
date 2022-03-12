package config

import (
	"raih-asa/beasiswa"
	"raih-asa/lomba"
	"raih-asa/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/raih_asa?parseTime=true"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&user.User{}, &beasiswa.Beasiswa{}, &beasiswa.CategoryBeasiswa{},
		&lomba.Lomba{}, &lomba.CategoryLomba{}, &lomba.Comment{}); err != nil {
		return nil, err
	}

	return db, nil
}
