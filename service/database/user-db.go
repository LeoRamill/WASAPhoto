package database

import (
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// Create a User --> Database Fuction that adds a new user in the database after registration
func (db *appdbimpl) CreateUser(user components.User) (string, error) {
	_, err := db.c.Exec(`
		INSERT INTO users (user_ID , username)
		VALUES (?,?)`,
		user.IdUser.Id, user.Usname)

	if err != nil {
		return user.IdUser.Id, err
	}

	return user.IdUser.Id, nil
}

// Get a Username by UserID
func (db *appdbimpl) GetUsername(id_user string) (string, error) {
	var username string
	err := db.c.QueryRow(`SELECT username FROM users WHERE user_ID LIKE ?`, id_user).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil

}

// Get a UserID by Username
func (db *appdbimpl) GetUserID(username components.Username) (id_user string, err error) {
	var ID string
	err = db.c.QueryRow(`
		SELECT user_ID
		FROM users
		WHERE username LIKE ?`,
		username.Usname).Scan(&ID)
	if err != nil {
		return "", err
	}
	return ID, nil

}

// Check a User --> Database Function that Checks if exists the user

func (db *appdbimpl) CheckUser(user components.User) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(user_ID)
		FROM users
		WHERE user_ID LIKE ? `,
		user.IdUser.Id).Scan(&count)
	// err dovrebbe essere un numero
	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return false, err
	}

	// se non c'è --> allora non esste
	if count == 0 {
		return false, err
	}

	// altrimenti esiste
	return true, nil
}

// Check a User --> Database Function that Checks if exists the user

func (db *appdbimpl) CheckUsername(usname string) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(username)
		FROM users
		WHERE username LIKE ? `,
		usname).Scan(&count)
	// err dovrebbe essere un numero
	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return false, err
	}

	// se non c'è --> allora non esste
	if count == 0 {
		return false, err
	}

	// altrimenti esiste
	return true, nil
}
