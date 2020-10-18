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

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", Login)
			auth.POST("/register", Register)
			auth.GET("/refresh/{rToken}", middlewares.AUTH(RefreshToken))
			auth.GET("/logout/{userId}", middlewares.AUTH(middlewares.TRUTH(Logout)))
		}

		users := v1.Group("/users")
		{
			users.GET("", GetUsers)
			users.GET("{userId}", middlewares.AUTH(middlewares.TRUTH(GetUser)))
			users.POST("", middlewares.AUTH(UpdateUser))
			users.DELETE("{userId}", middlewares.AUTH(middlewares.TRUTH(DeleteUser)))
		}
	}

	return router
}
