package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luka-sijic/flux/internal/models"

	"github.com/labstack/echo/v4"
)

func (h *FriendHandler) AddFriend(c echo.Context) error {
	username := c.Get("username").(string)
	user := new(models.FriendDTO)
	if err := c.Bind(&user); err != nil {
		log.Println(err)
	}

	fmt.Println(username)
	fmt.Println(user.Friend)

	result := h.svc.AddFriend(username, user)
	if !result {
		return c.JSON(http.StatusNotFound, "Failed to add friend")
	}

	return c.JSON(http.StatusOK, "Friend request sent")
}

func (h *FriendHandler) GetFriends(c echo.Context) error {
	id := c.Param("id")

	result := h.svc.GetFriends(id)
	if len(result) == 0 {
		return c.JSON(http.StatusNotFound, "No friends found")
	}
	return c.JSON(http.StatusOK, result)
}

func (h *FriendHandler) GetRequest(c echo.Context) error {
	id := c.Get("id").(string)

	result := h.svc.GetRequests(id)
	if len(result) == 0 {
		return c.JSON(http.StatusOK, "")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *FriendHandler) Respond(c echo.Context) error {
	username := c.Get("username").(string)
	action := new(models.FriendActionDTO)
	fmt.Println(username)
	if err := c.Bind(&action); err != nil {
		log.Println("Error?", err)
		return c.JSON(http.StatusInternalServerError, "error with friend data")
	}
	if action.FriendID == "" || action.Action == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing fields"})

	}
	result := h.svc.FriendResponse(username, action)
	if !result {
		fmt.Println("Error")
		return c.JSON(http.StatusInternalServerError, "Failed to update friend request")
	}
	return c.JSON(http.StatusOK, "Friend request updated")
}

func (h *FriendHandler) GetLog(c echo.Context) error {
	user1 := c.Param("user1")
	user2 := c.Param("user2")

	res := h.svc.GetLog(user1, user2)
	if res == nil {
		log.Panic("Error")
	}
	return c.JSON(http.StatusOK, res)
}
