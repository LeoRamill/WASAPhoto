package api

import (
	"encoding/json"

	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// /users/{user-id}/profile/photos/{photo-id}/comments/{comment-id}

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photo_id := ps.ByName("photo-id")
	user_id := ps.ByName("user-id")
	comment_id := ps.ByName("comment-id")
	// authorization
	id := r.Header.Get("Authorization")
	// check if exist path name
	if photo_id == "" || user_id == "" || comment_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting photo id or comment id")
		return
	}

	// Bad Request 400 --> length username
	if (len(user_id) < 6 && len(user_id) > 12) || (len(id) != 64) || (len(photo_id) != 64) || (len(comment_id) != 64) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	var imageId components.ImageID
	imageId.IDImage.Id = photo_id

	// get userID of photo author by imageId
	id_auth, err := rt.db.GetOwnerPhoto(imageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// Check if it exist the id_auth
	if id_auth == "" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// get username of photo author by id_usname
	photoAuthorUsername, err := rt.db.GetUsername(id_auth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// Check if it is a valid and exist photo
	if photoAuthorUsername == "" {
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
	usname_req, err := rt.db.GetUsername(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// Check if it is a valid and exist username
	if usname_req == "" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// CREATION VARIABLE
	var requestingUser components.User
	requestingUser.IdUser.Id = id
	requestingUser.Usname = usname_req

	var photoAuthorUser components.User
	photoAuthorUser.IdUser.Id = id_auth
	photoAuthorUser.Usname = photoAuthorUsername

	// 403
	banned_check, err := rt.db.CheckBanned(requestingUser, photoAuthorUser) // Fare Check banner
	if err != nil {
		ctx.Logger.WithError(err).Error("db.CheckBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	if banned_check {
		// User was banned by owner, can't post the comment
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(components.ForbiddenError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	var id_comment components.CommentID
	id_comment.IdComment.Id = comment_id
	comment_check, err := rt.db.CheckComment(id_comment)
	if err != nil {
		ctx.Logger.WithError(err).Error("db.CheckLike: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	if !comment_check {
		// you put like at most one time
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))
		ctx.Logger.WithError(err).Error("error checking comment existence, details 4")
		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// delete the comment in the database
	_, err = rt.db.UncommentPhoto(id_comment, requestingUser)
	if err != nil {
		// ctx.Logger.WithError(err).Error("put-like: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment: failed to execute query for insertion")
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)

}

// commentPhoto    /users/{user-id}/profile/photos/{photo-id}/comments/{comment-id}
func (rt *_router) createComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	photo_id := ps.ByName("photo-id")
	user_id := ps.ByName("user-id")

	// authorization
	id := r.Header.Get("Authorization")
	// check if exist path name
	if photo_id == "" || user_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting photo id or comment id")
		return
	}

	// Bad Request 400 --> length username
	if (len(user_id) < 6 && len(user_id) > 12) || (len(id) != 64) || (len(photo_id) != 64) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	var imageId components.ImageID
	imageId.IDImage.Id = photo_id

	// get userID of photo author by imageId
	id_auth, err := rt.db.GetOwnerPhoto(imageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	if id_auth == "" {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// get username of photo author by id_usname
	photoAuthorUsername, err := rt.db.GetUsername(id_auth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// Check if it is a valid and exist photo
	if photoAuthorUsername == "" {
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
	usname_req, err := rt.db.GetUsername(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// Check if it is a valid and exist username
	if usname_req == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// CREATION VARIABLE
	var requestingUser components.User
	requestingUser.IdUser.Id = id
	requestingUser.Usname = usname_req

	var photoAuthorUser components.User
	photoAuthorUser.IdUser.Id = id_auth
	photoAuthorUser.Usname = photoAuthorUsername

	// 403
	banned_check, err := rt.db.CheckBanned(requestingUser, photoAuthorUser) // Fare Check banner
	if err != nil {
		ctx.Logger.WithError(err).Error("db.CheckBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	if banned_check {
		// User was banned by owner, can't post the comment
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(components.ForbiddenError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// Retrieved comment from request body
	var comment components.Comment
	comment.DateTime = time.Now().UTC().String()
	comment.IdComment.IdComment.Id = RandomString(64)
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment/Decode: failed to decode request body json")
		return
	}

	// Check if the comment has a valid lenght
	if len(comment.Text) > 1024 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment: comment longer than 1024 characters")
		return
	}

	// Check if IDuser from request body matches to id requesting
	if comment.User.IdUser.Id != id {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment: comment longer than 1024 characters")
		return
	}

	// Check if IDphoto from request body matches to id photo
	if comment.IdPhoto.IDImage.Id != photo_id {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment/Decode: failed to decode request body json")
		return
	}

	// put/insert the comment in the database
	_, err = rt.db.SetPhotoComment(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment/db.CommentPhoto: failed to execute query for insertion")
		return
	}

	w.WriteHeader(http.StatusCreated)
	// The response body will contain the unique id of the comment
	err = json.NewEncoder(w).Encode(comment)
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
