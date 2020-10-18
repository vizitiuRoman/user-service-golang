package server

import (
	. "github.com/fasthttp/router"
	. "github.com/user-service/pkg/controllers"
	"github.com/user-service/pkg/middlewares"
)

func initControllers() *Router {
	router := New()

	// Ping
	router.GET("/api/ping", Ping)

	// Auth
	router.POST("/api/v1/auth/login", Login)
	router.POST("/api/v1/auth/register", Register)
	router.GET("/api/v1/auth/refresh/{rToken}", middlewares.AUTH(RefreshToken))
	router.GET("/api/v1/auth/logout/{userId}", middlewares.AUTH(middlewares.TRUTH(Logout)))

	// Users
	router.GET("/api/v1/users", GetUsers)
	router.GET("/api/v1/users/{userId}", middlewares.AUTH(middlewares.TRUTH(GetUser)))
	router.POST("/api/v1/users", middlewares.AUTH(UpdateUser))
	router.DELETE("/api/v1/users/{userId}", middlewares.AUTH(middlewares.TRUTH(DeleteUser)))

	return router
}
