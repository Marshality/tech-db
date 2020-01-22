package models

type Thread struct {
	ID uint64 `json:"id"`
	Slug string `json:"slug"`
	Author string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Title string `json:"title"`
	Votes uint64 `json:"votes"`
	CreatedAt string `json:"created"`
}
