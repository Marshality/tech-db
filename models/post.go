package models

type Post struct {
	ID        uint64 `json:"id"`
	Forum     string `json:"forum"`
	Thread    uint64 `json:"thread"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	Parent    uint64 `json:"parent"`
	IsEdited  bool   `json:"isEdited"`
	CreatedAt string `json:"created"`
}

type Posts []*Post
