package database

import (
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

/*
We have the following functions:
	0. CreatePhoto
	1. GetPhoto --> get photo information
	2. UpdatePhotoDescrp --> Update description of photo
	3. UpdatePhotoImg
	4. DeletePhoto --> Delete the photo
	5. CheckPhotoExists --> Check if the Photo Exist
	6. GetOwnerPhoto --> Get Owner of the Photo
*/

func (db *appdbimpl) CreatePostedPhoto(p components.PostedPhoto) (id string, err error) {
	_, err = db.c.Exec(`INSERT INTO photos (image_ID, poster_ID, descrp, dateTime, url)
						  VALUES (?,?,?,?,?)`,
		p.IdPhoto.Id, p.IdUser.Id, p.Descrp, p.DateTime, p.Url)
	if err != nil {
		return "", err
	}
	return p.IdPhoto.Id, err
}

// 1
func (db *appdbimpl) GetPhoto(targetPhoto components.ImageID, requestingUser components.User) (ph components.PostedPhoto, err error) {
	var photo components.PostedPhoto
	err = db.c.QueryRow(`SELECT *
						 FROM photos
						 WHERE (image_ID = ?) AND poster_ID NOT IN (SELECT banisher_ID FROM bans WHERE banned_ID = ? )`,
		targetPhoto.IDImage.Id, requestingUser.IdUser.Id).Scan(&photo.IdPhoto.Id, &photo.IdUser.Id, &photo.Descrp, &photo.DateTime, &photo.Url)
	if err != nil {
		return photo, err
	}
	photo.ListLike, err = db.GetPhotoLike(targetPhoto, requestingUser)
	if err != nil {
		return photo, err
	}

	photo.ListComment, err = db.GetPhotoComment(targetPhoto, requestingUser)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

// 2
func (db *appdbimpl) UpdatePhotoDescrp(targetPhoto components.ImageID, owner components.User, newDescrp string) (errstring string, err error) {
	_, err = db.c.Exec(
		`UPDATE photos 
	 	SET decrp LIKE ?
	 	WHERE (poster_ID LIKE ?) AND (image_ID LIKE ?)`,
		newDescrp, owner.IdUser.Id, targetPhoto.IDImage.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "done", nil
}

// 3
func (db *appdbimpl) UpdatePhotoImg(targetPhoto components.ImageID, owner components.User, newURL string) (errstring string, err error) {
	_, err = db.c.Exec(
		`UPDATE photos 
	 	SET decrp LIKE ?
	 	WHERE (poster_ID LIKE ?) AND (image_ID LIKE ?)`,
		newURL, owner.IdUser.Id, targetPhoto.IDImage.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "done", nil
}

// 4
func (db *appdbimpl) RemovePhoto(targetPhoto components.ImageID, owner components.User) (errstring string, err error) {
	_, err = db.c.Exec(
		`DELETE FROM photos 
	 	WHERE (poster_ID LIKE ?) AND (image_ID LIKE ?)`,
		owner.IdUser.Id, targetPhoto.IDImage.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "", nil
}

// 5

func (db *appdbimpl) CheckPhoto(photo components.ImageID) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(image_ID)
		FROM photos
		WHERE image_ID LIKE ? `,
		photo.IDImage.Id).Scan(&count)
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

// Get photo owner by photo_id
func (db *appdbimpl) GetOwnerPhoto(ph_ID components.ImageID) (user_id string, err error) {
	var user string
	err = db.c.QueryRow(`SELECT poster_ID
						 FROM photos
						 WHERE (image_ID = ?)`,
		ph_ID.IDImage.Id).Scan(&user)
	if err != nil {
		return "", err
	}
	return user, nil
}
