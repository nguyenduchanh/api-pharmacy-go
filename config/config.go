package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Các biến cấu hình MySQL
var (
	DBUser     = "vnpt_pharmacy_usr"
	DBPassword = "*%654Fo39S@1t$fga2e#sDk24A@"
	DBHost     = "10.30.29.14"
	DBPort     = "3306"
	DBName     = "go_pharmacy"
)

func ConnectDatabase() {
	// Tạo DSN từ các biến trên
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	// Kết nối MySQL qua GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối MySQL: %v", err)
	}

	DB = db
	fmt.Println("✅ Kết nối MySQL thành công")
}
