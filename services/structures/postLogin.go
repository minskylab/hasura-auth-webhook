package structures

// LOGIN types
type PostLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostLoginRes struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken"`
}
