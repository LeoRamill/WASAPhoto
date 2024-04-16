package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// Function that serves the requested photo

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	http.ServeFile(w, r,
		filepath.Join(photoFolder, ps.ByName("user-id"), "photos", ps.ByName("photo-id")))

}

func (rt *_router) deletePostedPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photo_id := ps.ByName("photo-id")
	photoAuthorUsername := ps.ByName("user-id")

	// Get the id by Authorization field
	id := r.Header.Get("Authorization")
	// check if exist the photo id or username
	if photo_id == "" || photoAuthorUsername == "" {
		// Doesn't exist the photo_id so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting photo id")

		return
	}

	// Bad Request 400 --> length author username, photo id
	if (len(photoAuthorUsername) < 6 && len(photoAuthorUsername) > 12) || (len(id) != 64) || (len(photo_id) != 64) {
		// BAD REQUEST: author username, photo id doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// get userID by photoauthorusername
	id_usname, err := rt.db.GetUserID(components.Username{Usname: photoAuthorUsername})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// check if exist the id
	if id_usname == "" {
		// STATUS NOT FOUND: id doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// get username by id requesting
	// var id_req components.UserID
	// id_req.IdUser.Id = id

	// get username by id
	usname_req, err := rt.db.GetUsername(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// check if exist the username
	if usname_req == "" {
		// STATUS NOT FOUND: id doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// responses 401: Unauthotorized
	// check if the id of the authorization field is the same of the
	// id searched with the username in the database
	if id != id_usname {
		// NO AUTHENTICATION: the two id don't coincide, so the action to search the username is unauthorized
		// write in the header that the status is unauthorized and return it
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.Error("error authenticating")
		return
	}

	// CREATION VARIABLE
	var photoAuthorUser components.User
	photoAuthorUser.IdUser.Id = id_usname
	photoAuthorUser.Usname = photoAuthorUsername

	var photoID components.ImageID
	photoID.IDImage.Id = photo_id

	// Call to the db function to remove the photo
	_, err = rt.db.RemovePhoto(photoID, photoAuthorUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-delete/RemovePhoto: error coming from database")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// Get the folder of the file that has to be eliminated
	pathPhoto, err := getUserPhotoFolder(photoAuthorUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-delete/getUserPhotoFolder: error with directories")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove the file from the user's photos folder
	err = os.Remove(filepath.Join(pathPhoto, photo_id))
	if err != nil {
		// Error occurs if the file doesn't exist, but for idempotency an error won't be raised
		ctx.Logger.WithError(err).Error("photo-delete/os.Remove: photo to be removed is missing")
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}

// users/{user-id}/profile/photos
func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// get username
	usname := ps.ByName("user-id")
	// get requesting  (UserID)
	id_req := r.Header.Get("Authorization")

	// 400 http status
	// check if exist the id
	if id_req == "" {
		// Doesn't exist the id_req so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting req id")
		return
	}

	// get userID by usname
	id_usname, err := rt.db.GetUserID(components.Username{Usname: usname})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// check if exist the id
	if id_usname == "" {
		// STATUS NOT FOUND: id doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// responses 401: Unauthotorized
	// check if the id of the authorization field is the same of the
	// id searched with the username in the database
	if id_req != id_usname {
		// NO AUTHENTICATION: the two id don't coincide, so the action to search the username is unauthorized
		// write in the header that the status is unauthorized and return it
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.Error("error authenticating")
		return
	}

	// Create a copy of the body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// After reading the body we won't be able to read it again. We'll re-assign a "fresh" io.ReadCloser to the body
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Check if the body content is either a png or a jpeg image, function implemented in utils.go
	err = checkFormatPhoto(r.Body, io.NopCloser(bytes.NewBuffer(data)), ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("photo-upload: body contains file that is neither jpg or png")
		// controllaerrore
		_ = json.NewEncoder(w).Encode("Image Format Error")
		return
	}
	// Body has been read in the previous function so it's necessary to reassign a io.ReadCloser to it
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Initialize photo struct
	var photo components.PostedPhoto
	photo.DateTime = time.Now().UTC().String()
	// function implemented in utils.go
	photo.IdPhoto.Id = RandomString(64)
	photo.IdUser.Id = id_usname
	photo.Usname = usname

	// Create the user's folder locally to save his/her images
	// implementation in utils.go
	PhotoPath, err := getUserPhotoFolder(usname)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error getting user's photo folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Create an empty file for storing the body content (image)
	out, err := os.Create(filepath.Join(PhotoPath, photo.IdPhoto.Id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error creating local photo file")
		return
	}

	// Copy body content to the previously created file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error copying body content into file photo")
		return
	}

	// Close the created file
	out.Close()

	photo.Url = PhotoPath

	// Create the post, so insert/put in the database the variable photo with type PostedPhoto
	_, err = rt.db.CreatePostedPhoto(photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error uploading photo")
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	// Send the output to the user
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment/Encode: failed convert photo_id to int64")
		return
	}
}
