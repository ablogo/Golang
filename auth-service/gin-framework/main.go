package main

import (
	"net/http"

	"gin-framework/api"
	"gin-framework/api/admin"
	"gin-framework/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//router.Use(middlewares.JWTMiddleware())
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/", index)

	{
		auth_router := router.Group("/auth")
		auth_router.POST("/sign-up", api.SignUp)
		auth_router.POST("/sign-in", api.SignIn)
		auth_router.POST("/validate-token", middlewares.JWTMiddleware(), api.VerifyToken)
	}
	{
		user_router := router.Group("/user", middlewares.JWTMiddleware())
		user_router.GET("/", api.GetUser)
		user_router.PUT("/user", api.UpdateUser)
		user_router.POST("/change-password", api.ChangePassword)
		user_router.GET("/img", api.GetPicture)
		user_router.POST("/img", api.AddPicture)
		user_router.GET("/address", api.GetAddresses)
		user_router.POST("/address", api.AddAddress)
		user_router.DELETE("/address", api.DeleteAddress)
		user_router.PUT("/address", api.UpdateAddress)
		user_router.DELETE("/user", api.DeleteUser)
	}
	{
		admin_router := router.Group("/admin", middlewares.JWTMiddleware())
		admin_router.GET("/user", admin.GetUser)
		admin_router.GET("/users", admin.GetUsers)
		admin_router.DELETE("/user", admin.DeleteUser)
		admin_router.GET("/user/address", admin.GetAddresses)
	}

	router.Run("localhost:8000")
}

func index(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "")
}
