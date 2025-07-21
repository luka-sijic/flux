package routes

import (
	"net/http"

	"github.com/luka-sijic/flux/internal/handler"
	"github.com/luka-sijic/flux/internal/service"
	"github.com/luka-sijic/flux/pkg/secret"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, svc *service.Infra) {
	userHandler := handler.NewUserHandler(svc)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	e.GET("/me", handler.Me)
	e.GET("/refresh", handler.Refresh)
}

func FriendRoutes(e *echo.Echo, svc *service.Infra) {
	friendHandler := handler.NewFriendHandler(svc)

	e.POST("/friend", friendHandler.AddFriend, secret.Auth)
	e.GET("/friend", friendHandler.GetRequest, secret.Auth)
	e.POST("/friend/respond", friendHandler.Respond, secret.Auth)
	e.GET("/friend/:id", friendHandler.GetFriends, secret.Auth)
	e.GET("/friend/:user1/:user2", friendHandler.GetLog)
}
