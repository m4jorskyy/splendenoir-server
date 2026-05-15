package data

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Login     string `json:"username"`
	Password  string `json:"password"`
}
