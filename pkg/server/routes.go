package server

import (
	. "github.com/fasthttp/router"
	. "github.com/user-service/pkg/controllers"
	"github.com/user-service/pkg/middlewares"
	"go.uber.org/zap"
)

func initControllers(logger *zap.SugaredLogger) *Router {
	router := New()

	c := NewController(logger)

	// Ping
	router.GET("/api/ping", c.Ping)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", c.Login)
			auth.POST("/register", c.Register)
			auth.GET("/refresh/{rToken}", middlewares.AUTH(c.RefreshToken))
			auth.GET("/logout/{userId}", middlewares.AUTH(middlewares.TRUTH(c.Logout)))
		}

		users := v1.Group("/users")
		{
			users.GET("", c.GetUsers)
			users.GET("{userId}", middlewares.AUTH(middlewares.TRUTH(c.GetUser)))
			users.POST("", middlewares.AUTH(c.UpdateUser))
			users.DELETE("{userId}", middlewares.AUTH(middlewares.TRUTH(c.DeleteUser)))
		}
	}

	return router
}
