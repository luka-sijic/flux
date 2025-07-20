package models

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FriendDTO struct {
	Friend string `json:"friend"`
}

type FriendActionDTO struct {
	FriendID string `json:"friendId"`
	Action   string `json:"action"`
}
