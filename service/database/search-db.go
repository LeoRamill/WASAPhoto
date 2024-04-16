package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// Get the user that I'm searching
func (db *appdbimpl) SearchUser(userToSearch string) ([]components.User, error) {
	// Prendiamo le colonne in cui c'è l'id e lo username del utente che stiamo cercando e allo stesso tempo
	// ci assicuriamo che non siamo stati bannati (che il searcher non è stato bannato)
	rows, err := db.c.Query(
		`SELECT * 
		 FROM users 
		 WHERE (username LIKE ?)`,
		userToSearch)

	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				logrus.Errorf("error closing result set: %v", err)
			}
		}
	}()

	if err != nil {
		return nil, fmt.Errorf("error searching user: %w", err)
	}

	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset.
	var remaining_user []components.User

	for rows.Next() {

		// Check Error
		if rows.Err() != nil {
			return nil, fmt.Errorf("error getting next user: %w", rows.Err())
		}

		// user := components.User{}
		var user components.User
		err = rows.Scan(&user.IdUser.Id, &user.Usname)
		if err != nil {
			return nil, err
		}
		// otherwise err == nil --> there isn't error
		remaining_user = append(remaining_user, user)
	}

	return remaining_user, nil
}
