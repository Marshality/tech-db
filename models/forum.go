package models

type Forum struct {
	ID      uint64 `json:"id,omitempty"`
	Posts   uint64 `json:"posts,omitempty"`
	Slug    string `json:"slug"`
	Threads uint64 `json:"threads,omitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}
