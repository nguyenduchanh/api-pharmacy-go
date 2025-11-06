package routes

import (
	"api-pharmacy-go/controllers"
	_ "api-pharmacy-go/docs"
	"api-pharmacy-go/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
)

func SetupRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RateLimitMiddleware(30, time.Minute))
	router.Use(middleware.SecurityHeaders())
	auth := router.Group("/api/auth")
	{
		//authen
		auth.GET("/register", controllers.Register)
		auth.GET("/login", controllers.Login)

	}
	users := router.Group("/api/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/", controllers.GetUsers)
		users.GET("/:id", controllers.GetUser)
		users.POST("/", controllers.CreateUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
		users.GET("/export", controllers.ExportUsers)
		users.POST("/import", controllers.ImportUsers)
	}
	role := router.Group("/api/roles")
	role.Use(middleware.AuthMiddleware())
	{
		role.GET("/", controllers.GetRoles)
		role.GET("/:id", controllers.GetRole)
		role.POST("/", controllers.CreateRole)
		role.PUT("/:id", controllers.UpdateRole)
		role.DELETE("/:id", controllers.DeleteRole)
	}
	emp := router.Group("/api/emps")
	emp.Use(middleware.AuthMiddleware())
	{
		emp.GET("/", controllers.GetEmps)
		emp.GET("/:id", controllers.GetEmp)
		emp.POST("/", controllers.CreateEmp)
		emp.PUT("/:id", controllers.UpdateEmp)
		emp.DELETE("/:id", controllers.DeleteEmp)
	}
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	er := router.Run(":8888")
	if er != nil {
		log.Fatal("Lỗi khởi động server:", er)
	}
}
