package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

/*
Functions:
	1. GetComment --> Get all information of comment
	2. SetComment --> Create  comment
	3. DeleteComment --> Delete comment
	4. UpdateDescrpComment --> Update Comment
	5. CheckComment --> Check if Exist Comment
*/

// 1

func (db *appdbimpl) GetPhotoComment(targetPhoto components.ImageID, requestingUser components.User) ([]components.Comment, error) {
	rows, err := db.c.Query(
		`SELECT * 
		 FROM comments
		 WHERE (photo_ID LIKE ?) AND (commenter_ID NOT IN (SELECT banisher_ID FROM bans WHERE banned_ID = ? ))`,
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
	var commentList []components.Comment

	for rows.Next() {

		// Check Error
		if rows.Err() != nil {
			return nil, fmt.Errorf("error getting next user: %w", rows.Err())
		}

		var comm components.Comment
		// allocazione delle variabili con i puntatori
		// Vedere se persiste ancora il problema
		err = rows.Scan(&comm.IdComment.IdComment.Id, &comm.IdPhoto.IDImage.Id, &comm.Text, &comm.DateTime, &comm.User.IdUser.Id)
		if err != nil {
			return nil, err
		}

		usname, err := db.GetUsername(comm.User.IdUser.Id)
		if err != nil {
			return nil, fmt.Errorf("error comment PostedPhoto: %w", err)
		}

		comm.Nickname = usname

		// If it's different with nil --> there is error
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		// otherwise err == nil --> there isn't error
		commentList = append(commentList, comm)
	}

	return commentList, nil
}

// 2
func (db *appdbimpl) SetPhotoComment(c components.Comment) (errstring string, err error) {
	_, err = db.c.Exec("INSERT INTO comments (comment_ID, photo_ID , text, dateTime, commenter_ID) VALUES (?,?,?,?,?)",
		c.IdComment.IdComment.Id, c.IdPhoto.IDImage.Id, c.Text, c.DateTime, c.User.IdUser.Id)
	if err != nil {
		return components.InternalServerError, fmt.Errorf("error changing username: %w", err)
	}
	return "", nil
}

// 3
func (db *appdbimpl) UncommentPhoto(targetComment components.CommentID, owner components.User) (errstring string, err error) {
	_, err = db.c.Exec(
		`DELETE FROM comments
	 	WHERE (commenter_ID LIKE ?) AND (comment_ID LIKE ?)`,
		owner.IdUser.Id, targetComment.IdComment.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "", nil
}

// 4
func (db *appdbimpl) UpdateCommentDescrp(targetComment components.CommentID, owner components.User, newText string) (errstring string, err error) {
	_, err = db.c.Exec(
		`UPDATE comments 
	 	SET text LIKE ?
	 	WHERE (commenter_ID LIKE ?) AND (comment_ID LIKE ?)`,
		newText, owner.IdUser.Id, targetComment.IdComment.Id)

	if err != nil {
		return components.InternalServerError, err
	}

	return "", nil
}

// 5
func (db *appdbimpl) CheckComment(comm components.CommentID) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(comment_ID)
		FROM comments
		WHERE comment_ID LIKE ? `,
		comm.IdComment.Id).Scan(&count)
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
