package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// /users/{user-id}/profile/photos/{photo-id}/likes/{like-id}
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// w = object ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	/*
		photo_id := ps.ByName("photo-id"): This line is likely part of a router handling an HTTP request. It looks like it's trying to retrieve a parameter named "photo-id" from the URL path. The value of this parameter is then being assigned to the variable photo_id. ps is likely an instance of httprouter.Params, which is a collection of route parameters parsed from the URL.
		user_id := ps.ByName("user-id"): Similar to the previous line, this is retrieving a parameter named "user-id" from the URL path and assigning its value to the variable user_id.
		like_id := ps.ByName("like-id"): Again, this line is retrieving a parameter named "like-id" from the URL path and assigning its value to the variable like_id.
	*/
	photo_id := ps.ByName("photo-id")
	user_id := ps.ByName("user-id")
	like_id := ps.ByName("like-id")

	/* AUTHORIZATION
	id will contain the value of the "Authorization" header, which often holds some form of authentication token or credentials.
	*/

	// r = object Request
	id := r.Header.Get("Authorization")
	// check if exist path name
	if photo_id == "" || user_id == "" || like_id == "" {
		// Doesn't exist the  photo/user/like id so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting photo id or comment id")
		return
	}

	// Bad Request 400 --> length username
	if (len(user_id) < 6 && len(user_id) > 12) || (len(id) != 64) || (len(photo_id) != 64) || (len(like_id) != 64) {
		// BAD REQUEST: username or id_req doesn't satisfy the costraints setting to the api documentation
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

	// get username of the request
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
		// User was banned by owner, can't put like
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(components.ForbiddenError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// var id_like components.LikeID
	// id_like.IdLike.Id = like_id

	var like_struct components.Like
	err = json.NewDecoder(r.Body).Decode(&like_struct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment/Decode: failed to decode request body json")
		return
	}

	// Check if IDlike from request body matches to like_id
	if like_struct.IdLike.IdLike.Id != like_id {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment: comment longer than 1024 characters")
		return
	}

	// Check if IDphoto from request body matches to photo_id
	if like_struct.IdPhoto.IDImage.Id != photo_id {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment: comment longer than 1025 characters")
		return
	}

	// Check if IDuser from request body matches to user_id
	if like_struct.User != usname_req {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("post-comment: comment longer than 1026 characters")
		return
	}

	// Insert the like in the db via db function
	_, err = rt.db.SetPhotoLike(like_struct)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-like: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// Respond with 201 http status
	w.WriteHeader(http.StatusCreated)
	// Send the output to the user. Instead of giving null for no matches return and empty slice of photos. ( ontrollaerrore)
	_ = json.NewEncoder(w).Encode(like_struct)

}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photo_id := ps.ByName("photo-id")
	user_id := ps.ByName("user-id")
	like_id := ps.ByName("like-id")
	// authorization
	id := r.Header.Get("Authorization")
	// check if exist path name
	if photo_id == "" || user_id == "" || like_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting photo id or comment id")
		return
	}

	// Bad Request 400 --> length username
	if (len(user_id) < 6 && len(user_id) > 12) || (len(id) != 64) || (len(photo_id) != 64) || (len(like_id) != 64) {
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
	// Check if it exist the id_auth --> DA RIVEDERE (forse eliminare)
	if id_auth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(components.UnauthorizedError))

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
	// No logged
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

	// 403
	var id_like components.LikeID
	id_like.IdLike.Id = like_id
	like_check, err := rt.db.CheckLike(id_like)
	if err != nil {
		ctx.Logger.WithError(err).Error("db.CheckLike: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	if !like_check {
		// you put like at most one time
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// delete the like in the database
	_, err = rt.db.UnLikePhoto(id_like, requestingUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)

}
