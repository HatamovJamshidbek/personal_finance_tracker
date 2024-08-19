package api

import (
	"auth_service/api/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"time"
)

// New ...
// @title           MyRentUz API
// @version         1
// @description     My Rent Uz
// @in header
// @name Authorization
func New(h handler.Handler) *gin.Engine {
	router := gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		RequestHeaders: "Authorization,Origins,Content-Type",
		Methods:        "POST, GET, PUT, DELETE, OPTIONS",
	}))
	auth := router.Group("/auth_service")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		//auth.PUT("/reset_password", h.ResetPassword)
		auth.GET("/refresh", h.Refresh)
		auth.GET("/logout", h.Logout)

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router

}
func traceRequest(c *gin.Context) {
	beforeRequest(c)

	c.Next()

	afterRequest(c)
}

func beforeRequest(c *gin.Context) {
	startTime := time.Now()
	c.Set("start_time", startTime)
	log.Println("start time:", startTime.Format("2006-01-02 15:04:05.0000"), "path:", c.Request.URL.Path)
}
func afterRequest(c *gin.Context) {

	startTime, exists := c.Get("start_time")
	if !exists {
		startTime = time.Now()
	}

	duration := time.Since(startTime.(time.Time)).Seconds()

	log.Println("end time:", time.Now().Format("2006-01-02 15:04:05.0000"), "duration:", duration, " second", "method:", c.Request.Method)
	fmt.Println()
}
