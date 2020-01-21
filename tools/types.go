package tools

type Body map[string]interface{}

type Error struct {
	Message string `json:"message"`
}
