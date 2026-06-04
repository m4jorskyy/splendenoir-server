package repositories

import (
	"database/sql"
	"splendenoir-server/internal/models/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(profile *user.UserProfile, creds *user.UserCredentials) (int64, error) {
	tx, errTx := r.db.Begin()
	if errTx != nil {
		return 0, errTx
	}

	defer func(tx *sql.Tx) {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return
		}
	}(tx)

	var returnedID int64

	errReturnID := tx.QueryRow("INSERT INTO user_profiles (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING user_profile.id", profile.FirstName, profile.LastName, profile.Email).Scan(&returnedID)

	if errReturnID != nil {
		return 0, errReturnID
	}

	_, errExec := tx.Exec("INSERT INTO user_credentials (username, password_hash, profile_id) VALUES ($1, $2, $3)", creds.Username, creds.PasswordHash, returnedID)

	if errExec != nil {
		return 0, errExec
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return 0, errCommit
	}

	return returnedID, nil
}

func (r *UserRepository) LoginUser(username string) (int64, string, string, error) {
	var returnedID int64
	var returnedPasswordHash string
	var returnedFirstName string

	errReturned := r.db.QueryRow("SELECT user_credentials.profile_id, user_credentials.password_hash, user_profiles.first_name FROM user_credentials, user_profiles WHERE username = $1 AND user_credentials.profile_id = user_profiles.id", username).Scan(&returnedID, &returnedPasswordHash, &returnedFirstName)

	if errReturned != nil {
		return 0, "", "", errReturned
	}

	return returnedID, returnedPasswordHash, returnedFirstName, nil

}
