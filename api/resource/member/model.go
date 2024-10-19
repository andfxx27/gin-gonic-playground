package member

import (
	"encoding/json"
	"time"
)

type Member struct {
	ID        string
	Username  string
	Password  string `json:"-"`
	Email     string
	CreatedAt time.Time
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required,alphanum,max=50"`
	Password string `json:"password" binding:"required,max=32"`
	Email    string `json:"email" binding:"required,email,max=255"`
}

type SignInRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func (m Member) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
