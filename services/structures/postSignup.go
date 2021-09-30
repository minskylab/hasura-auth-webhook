package structures

// SIGNUP types
type PostSignupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostSignupRes struct {
	UserID string `json:"userId"`
}
