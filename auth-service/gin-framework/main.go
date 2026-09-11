package main

import (
	"net/http"
	database "src/db"
	"src/services"
	"src/utils"

	"gin-framework/api"
	"gin-framework/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	settings := &utils.Settings{}
	settings.GetInstance()
	postgresHandler := &database.DBHandler{Settings: settings}
	userSvc := &services.UserService{Postgres: postgresHandler}
	routerHandler := &api.Router{UserSvc: userSvc}
	adminRouterHandler := &api.AdminRouter{UserSvc: userSvc}

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
		auth_router.POST("/sign-up", routerHandler.SignUp)
		auth_router.POST("/sign-in", routerHandler.SignIn)
		auth_router.POST("/validate-token", middlewares.JWTMiddleware(), api.VerifyToken)
	}
	{
		user_router := router.Group("/user", middlewares.JWTMiddleware())
		user_router.GET("/", routerHandler.GetUser)
		user_router.PUT("/user", routerHandler.UpdateUser)
		user_router.POST("/change-password", routerHandler.ChangePassword)
		user_router.GET("/img", routerHandler.GetPicture)
		user_router.POST("/img", routerHandler.AddPicture)
		user_router.GET("/address", routerHandler.GetAddresses)
		user_router.POST("/address", routerHandler.AddAddress)
		user_router.DELETE("/address", routerHandler.DeleteAddress)
		user_router.PUT("/address", routerHandler.UpdateAddress)
		user_router.DELETE("/user", routerHandler.DeleteUser)
	}
	{
		admin_router := router.Group("/admin", middlewares.JWTMiddleware())
		admin_router.GET("/user", adminRouterHandler.GetUser)
		admin_router.GET("/users", adminRouterHandler.GetUsers)
		admin_router.DELETE("/user", adminRouterHandler.DeleteUser)
		admin_router.GET("/user/address", adminRouterHandler.GetAddresses)
	}

	router.Run("localhost:8000")
}

func index(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "")
}
