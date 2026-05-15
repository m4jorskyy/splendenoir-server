package data

type LoginRequest struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}
