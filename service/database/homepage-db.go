package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// Function that gets the stream of user --> List of PostedPhoto
func (db *appdbimpl) GetStream(user components.User) ([]components.PostedPhoto, error) {
	rows, err := db.c.Query(
		`SELECT p.image_ID, p.poster_ID, p.descrp, p.dateTime, p.url 
		 FROM photos p, users u
		 WHERE (p.poster_ID = u.user_ID) AND (p.poster_ID IN (SELECT followed_ID FROM followers WHERE followed_ID LIKE ? ))
		 ORDER BY p.dateTime DESC`,
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
