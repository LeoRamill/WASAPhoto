package api

import (
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Create Random String
func RandomString(length int) string {
	return stringWithCharset(length, charset)
}

// Function that creates a new subdir for the specified user
func createUserFolder(identifier string, ctx reqcontext.RequestContext) error {

	// Create the path media/useridentifier/ inside the project dir
	path := filepath.Join(photoFolder, identifier)

	// To the previously created path add the "photos" subdir
	err := os.MkdirAll(filepath.Join(path, "photos"), os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("session/createUserFolder:: error creating directories for user")
		return err
	}
	return nil
}

// Function that returns the path of the photo folder for a certain user
func getUserPhotoFolder(user_id string) (UserPhotoFoldrPath string, err error) {

	// Path of the photo dir "./media/user_id/photos/"
	photoPath := filepath.Join(photoFolder, user_id, "photos")

	return photoPath, nil
}

// Function checks if the format of the photo is png or jpeg. Returns the format extension and an error
func checkFormatPhoto(body io.ReadCloser, newReader io.ReadCloser, ctx reqcontext.RequestContext) error {

	/*
	 jpeg.Decode() function to decode a JPEG image from a byte slice named body. This function returns two values: the decoded image and an error.
	*/
	_, errJpg := jpeg.Decode(body)
	// Check if exist an JPEG error
	if errJpg != nil {

		// same procedure above but with png image
		body = newReader
		_, errPng := png.Decode(body)
		if errPng != nil {
			return errPng
		}
		return nil
	}
	return nil
}
