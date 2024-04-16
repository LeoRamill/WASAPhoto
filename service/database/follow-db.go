package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

/*
Functions:
	1. GetFollowed --> Get List of following/followers user
	2. FollowUser --> Follow User
	3. UnFollowUser -->  Unfollow the User
	4. UpdateDescrpComment --> Update Comment
	5. CheckFollow --> Check if Exist the following
	6. GetFollowers

	Potremmo mettere getUserbyComment
*/

// Get list of followed by userID
func (db *appdbimpl) GetFollowed(user components.User) (followeds []components.User, err error) {
	// Estrazione dei followed ID
	rows, err := db.c.Query(
		`SELECT f.followed_ID, u.username
		 FROM users u, followers f
		 WHERE (f.followed_ID = u.user_ID) AND (follower_ID LIKE ?) `,
		user.IdUser.Id)

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

	var list_followed []components.User
	for rows.Next() {
		if rows.Err() != nil {
			return nil, fmt.Errorf("error getting next user: %w", rows.Err())
		}

		var followed components.User
		err = rows.Scan(&followed.IdUser.Id, &followed.Usname)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		list_followed = append(list_followed, followed)
	}
	return list_followed, nil
}

// 2
func (db *appdbimpl) FollowUser(follower components.User, followed components.User) (errstring string, err error) {
	_, err = db.c.Exec("INSERT INTO followers (follower_ID, followed_ID) VALUES (?,?)",
		follower.IdUser.Id, followed.IdUser.Id)
	if err != nil {
		return components.InternalServerError, fmt.Errorf("error changing userID: %w", err)
	}
	return "", nil
}

// 3
func (db *appdbimpl) UnFollowUser(follower components.User, followed components.User) (errstring string, err error) {
	_, err = db.c.Exec(
		`DELETE FROM followers
	 	WHERE (follower_ID LIKE ?) AND (followed_ID LIKE ?)`,
		follower.IdUser.Id, followed.IdUser.Id)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error changing username: %w", err)
	}

	return "", nil
}

// 4
func (db *appdbimpl) CheckFollow(follower components.User, followed components.User) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(followed_ID)
		FROM followers
		WHERE (followed_ID LIKE ?) AND (follower_ID LIKE ?) `,
		followed.IdUser.Id, follower.IdUser.Id).Scan(&count)
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

// Get list of followed by userID
func (db *appdbimpl) GetFollowers(user components.User) (followeds []components.User, err error) {
	// Estrazione dei followed ID
	rows, err := db.c.Query(
		`SELECT f.follower_ID, u.username
		 FROM users u, followers f
		 WHERE (f.follower_ID = u.user_ID) AND (followed_ID LIKE ?) `,
		user.IdUser.Id)

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

	var list_followed []components.User
	for rows.Next() {
		if rows.Err() != nil {
			return nil, fmt.Errorf("error getting next user: %w", rows.Err())
		}

		var followed components.User
		err = rows.Scan(&followed.IdUser.Id, &followed.Usname)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		list_followed = append(list_followed, followed)
	}
	return list_followed, nil
}
