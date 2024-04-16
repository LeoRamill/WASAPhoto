package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

/*
Create the following function:
	1. GetBans --> Get List of followed by userID
	2. BanUser --> Essentially ban the user
	3. UnBanUser --> Essentially remove the ban of user
	4. CheckBanned --> Check if ban exist
*/

// 1.
func (db *appdbimpl) GetBans(user components.User) (bans []components.User, err error) {
	// Estrazione dei followed ID
	rows, err := db.c.Query(
		`SELECT b.banned_ID, u.username
		 FROM users u, bans b
		 WHERE (b.banned_ID = u.user_ID) AND (b.banisher_ID LIKE ?)`,
		user.IdUser.Id)

	/*
		Anonymous function: The defer statement is followed by an anonymous function.
		This function will be executed when the surrounding function returns.
	*/
	defer func() {
		if rows != nil {
			err := rows.Close()
			if err != nil {
				logrus.Errorf("error closing result set: %v", err)
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	var list_ban []components.User
	for rows.Next() {
		if rows.Err() != nil {
			return nil, fmt.Errorf("error getting next user: %w", rows.Err())
		}

		var ban components.User
		err = rows.Scan(&ban.IdUser.Id, &ban.Usname)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		list_ban = append(list_ban, ban)
	}
	return list_ban, nil
}

// 2.
func (db *appdbimpl) BanUser(banner components.User, banned components.User) (errstring string, err error) {
	// INSERT INTO THE BAN DATABASE THE BAN
	_, err = db.c.Exec("INSERT INTO bans (banisher_ID, banned_ID) VALUES (?,?)",
		banner.IdUser.Id, banned.IdUser.Id)
	if err != nil {
		return components.InternalServerError, err
	}
	// REMOVE INTO THE FOLLOWER DATABASE THE BANNER
	_, err = db.c.Exec(
		`DELETE FROM followers
	 	WHERE (follower_ID LIKE ?) AND (followed_ID LIKE ?)`,
		banner.IdUser.Id, banned.IdUser.Id)

	if err != nil {
		return components.InternalServerError, err
	}
	return "", nil
}

// 3
func (db *appdbimpl) UnBanUser(banner components.User, banned components.User) (errstring string, err error) {
	_, err = db.c.Exec(
		`DELETE FROM bans
	 	WHERE (banisher_ID LIKE ?) AND (banned_ID LIKE ?)`,
		banner.IdUser.Id, banned.IdUser.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "", nil
}

// 4
func (db *appdbimpl) CheckBanned(banish components.User, ban components.User) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*)
		FROM bans
		WHERE (banned_ID LIKE ?) AND (banisher_ID LIKE ?) `,
		ban.IdUser.Id, banish.IdUser.Id).Scan(&count)
	// err dovrebbe essere un numero
	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}

	// se non c'è --> allora non esste
	if count == 1 {
		return true, nil
	}

	// altrimenti esiste
	return false, nil
}
