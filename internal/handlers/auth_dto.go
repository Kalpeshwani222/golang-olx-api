package handlers

import (
	"strings"
	"time"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return ValidationError{Field: "email", Msg: "must not be empty"}
	}
	if !strings.Contains(r.Email, "@") {
		return ValidationError{Field: "email", Msg: "must be a valid email address"}
	}
	if len(r.Password) < 6 {
		return ValidationError{Field: "password", Msg: "must be at least 6 characters"}
	}
	return nil
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

//Login  ----

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return ValidationError{Field: "email", Msg: "must not be empty"}
	}
	if strings.TrimSpace(r.Password) == "" {
		return ValidationError{Field: "password", Msg: "must not be empty"}
	}
	return nil
}


type LoginResponse struct {
	Token string `json:"token"`
}


type User struct {
	ID        string
	Email     string
	Password  string    // bcrypt hash
	CreatedAt time.Time
}

