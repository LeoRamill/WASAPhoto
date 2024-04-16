package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

/*
Create two functions:
	1. Get Gallery of Personal Photos
	2. Update the Username
	3. Post Photo
	4. UncommentOwnerPhoto --> Delete every comment in owner photo
*/

func (db *appdbimpl) GetGallery(user components.User) (stream []components.PostedPhoto, err error) {
	rows, err := db.c.Query(
		`SELECT p.image_ID, p.poster_ID, p.descrp, p.dateTime, p.url 
		 FROM photos p, users u
		 WHERE (p.poster_ID = u.user_ID) AND (p.poster_ID LIKE ?)`,
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
		return nil, fmt.Errorf("error searching userID: %w", err)
	}

	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset.
	var photos []components.PostedPhoto

	for rows.Next() {
		// Check Error
		if rows.Err() != nil {
			return nil, fmt.Errorf("error getting next PostedPhoto: %w", rows.Err())
		}

		// photo := components.PostedPhoto{}
		var photo components.PostedPhoto
		// allocazione delle variabili con i puntatori
		err = rows.Scan(&photo.IdPhoto.Id, &photo.IdUser.Id, &photo.Descrp, &photo.DateTime, &photo.Url)
		// If it's different with nil --> there is error
		if err != nil {
			return nil, fmt.Errorf("error scanning PostedPhoto: %w", err)
		}

		usname, err := db.GetUsername(user.IdUser.Id)
		if err != nil {
			return nil, fmt.Errorf("error comment PostedPhoto: %w", err)
		}
		photo.Usname = usname

		// Get the list of comment
		comments, err := db.GetPhotoComment(components.ImageID{IDImage: photo.IdPhoto}, user)
		if err != nil {
			return nil, fmt.Errorf("error comment PostedPhoto: %w", err)
		}
		photo.ListComment = comments

		likes, err := db.GetPhotoLike(components.ImageID{IDImage: photo.IdPhoto}, user)
		if err != nil {
			return nil, fmt.Errorf("error like PostedPhoto: %w", err)
		}
		photo.ListLike = likes

		// otherwise err == nil --> there isn't error
		photos = append(photos, photo)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error rows PostedPhoto: %w", err)
	}

	return photos, nil

}

func (db *appdbimpl) UpdateUsername(user components.User, newUsername string) (errstring string, err error) {
	_, err = db.c.Exec(`UPDATE users SET username = ? WHERE (user_ID LIKE ?) AND (username LIKE ?)`, newUsername, user.IdUser.Id, user.Usname)

	if err != nil {
		return components.InternalServerError, err
	}

	return newUsername, nil
}

func (db *appdbimpl) PostPhoto(ph components.PostedPhoto) (errstring string, err error) {
	_, err = db.c.Exec(`INSERT INTO photos (image_ID, poster_ID , descrp, dateTime, url) VALUES (?,?,?,?,?)`,
		ph.IdPhoto.Id, ph.IdUser.Id, ph.Descrp, ph.DateTime, ph.Url)
	if err != nil {
		return components.InternalServerError, err
	}
	return "", nil
}

// funzione che può eliminare i commenti dalle proprie foto
func (db *appdbimpl) UncommentOwnerPhoto(targetComment components.CommentID) (errstring string, err error) {
	_, err = db.c.Exec(
		`DELETE FROM comments
	 	WHERE (comment_ID LIKE ?)`,
		targetComment.IdComment.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "", nil
}
