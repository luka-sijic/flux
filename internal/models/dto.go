package models

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type FriendDTO struct {
	Friend string `json:"friend"`
}

type FriendActionDTO struct {
	FriendID string `json:"friendId"`
	Action   string `json:"action"`
}
