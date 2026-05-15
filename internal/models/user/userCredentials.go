package user

type UserCredentials struct {
	ID           int64
	Username     string
	PasswordHash string
	ProfileID    int64
}
