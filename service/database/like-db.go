package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

/*
Functions:
	1. GetLike
	2. SetLike
	3. DeleteLike
*/

// 1
func (db *appdbimpl) GetPhotoLike(targetPhoto components.ImageID, requestingUser components.User) (likes []components.Like, err error) {
	rows, err := db.c.Query(
		`SELECT * 
		 FROM likes
		 WHERE (photo_ID LIKE ?) AND (like_ID NOT IN (SELECT banisher_ID FROM bans WHERE banned_ID = ? ))`,
		targetPhoto.IDImage.Id, requestingUser.IdUser.Id)

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

	// Read all the users in the resulset.
	var likeList []components.Like

	for rows.Next() {

		// Check Error
		if rows.Err() != nil {
			return nil, fmt.Errorf("error getting next like: %w", rows.Err())
		}

		var l components.Like
		// allocazione delle variabili con i puntatori
		err = rows.Scan(&l.IdLike.IdLike.Id, &l.IdPhoto.IDImage.Id, &l.User)
		// mancherebbero : ImageData , ListLike, ListComment --> vediamo come riuscirle a mettere

		// If it's different with nil --> there is error
		if err != nil {
			return nil, fmt.Errorf("error scanning like: %w", err)
		}
		// otherwise err == nil --> there isn't error
		likeList = append(likeList, l)
	}

	return likeList, nil

}

// 2
func (db *appdbimpl) SetPhotoLike(l components.Like) (errstring string, err error) {
	_, err = db.c.Exec("INSERT INTO likes (like_ID, photo_ID , liker_ID) VALUES (?,?,?)",
		l.IdLike.IdLike.Id, l.IdPhoto.IDImage.Id, l.User)
	if err != nil {
		return components.InternalServerError, err
	}
	return "", nil
}

// 3
func (db *appdbimpl) UnLikePhoto(l components.LikeID, owner components.User) (errstring string, err error) {
	_, err = db.c.Exec(
		`DELETE FROM likes
	 	WHERE (like_ID LIKE ?) AND (like_ID LIKE ?)`,
		owner.IdUser.Id, l.IdLike.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "", nil
}

func (db *appdbimpl) CheckLike(l components.LikeID) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(like_ID)
		FROM likes
		WHERE like_ID LIKE ? `,
		l.IdLike.Id).Scan(&count)
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
