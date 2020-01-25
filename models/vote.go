package models

type Vote struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   uint64 `json:"thread"`
}
