package server

import "time"

type NewRegisterEvent struct {
	EventAt time.Time `json:"event_at"`
	UserID  string    `json:"user_id"`
}

type LoginMagicLinkEvent struct {
	EventAt time.Time `json:"event_at"`
	UserID  string    `json:"user_id"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RecoverPasswordEvent struct {
	EventAt time.Time `json:"event_at"`
	UserID  string    `json:"user_id"`

	Token string `json:"token"`
}
