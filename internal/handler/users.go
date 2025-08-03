package handler

import (
	"log"
	"net/http"

	"github.com/luka-sijic/flux/internal/models"

	"github.com/luka-sijic/flux/pkg/secret"

	"github.com/labstack/echo/v4"
)

func (h *UserHandler) Register(c echo.Context) error {
	user := new(models.UserDTO)
	if err := c.Bind(user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to bind to user"})
	}

	result := h.svc.CreateUser(user)
	if !result {
		return c.JSON(http.StatusInternalServerError, "Failed to register user")
	}

	return c.JSON(http.StatusOK, "User successfully registered")
}

func (h *UserHandler) Login(c echo.Context) error {
	user := new(models.UserDTO)
	if err := c.Bind(&user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, "Error binding struct")
	}

	result := h.svc.LoginUser(user)
	if result == nil {
		return c.JSON(http.StatusInternalServerError, "Incorrect username/password")
	}

	access := secret.GenerateJWT(&models.User{ID: result.ID, Username: result.Username}, 900)
	refresh := secret.GenerateJWT(&models.User{ID: result.ID, Username: result.Username}, 2592000)

	setCookie(c, "access", access, 900)
	setCookie(c, "refresh", refresh, 2592000)

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) Profile(c echo.Context) error {
	username := c.Param("username")
	res := h.svc.Profile(username)
	if res {
		return c.JSON(http.StatusOK, "1")
	}
	return c.JSON(http.StatusOK, "0")
}
