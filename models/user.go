package models

type User struct {
	ID       uint64 `json:"id,omitempty"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	About    string `json:"about"`
}

func (u *User) Map(other *User) {
	if other.Email == "" {
		other.Email = u.Email
	}

	if other.Fullname == "" {
		other.Fullname = u.Fullname
	}

	if other.About == "" {
		other.About = u.About
	}
}
