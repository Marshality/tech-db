package models

type User struct {
	ID       uint64 `json:"id"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	About    string `json:"about"`
}
