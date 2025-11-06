package main

import (
	"api-pharmacy-go/config"
	_ "api-pharmacy-go/docs"
	"api-pharmacy-go/models"
	"api-pharmacy-go/routes"
	"fmt"
)

// @title Pharmacy API
// @version 1.0
// @description API for Pharmacy Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.ConnectDatabase()
	err := config.DB.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(
			&models.MUser{},
			models.MRole{},
			&models.UserRole{},
			&models.MEmp{})
	if err != nil {
		fmt.Println("Lỗi migrate:", err)
		return
	}
	//r := routes.SetupRouter()
	//r.Run(":8080")
	routes.SetupRouter()
}
