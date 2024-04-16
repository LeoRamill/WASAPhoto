package api

import (
	"encoding/json"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// /users/{user-id}/homepage:
func (rt *_router) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// Get the id by Authorization field
	id_req := r.Header.Get("Authorization")

	usname := ps.ByName("user-id")

	// check if exist the id or username
	if id_req == "" || usname == "" {
		// Doesn't exist the username or id_req so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting requesting id or username1")
		return
	}

	// Bad Request 400 --> length username or id
	if (len(usname) < 6 && len(usname) > 12) || (len(id_req) != 64) {
		// BAD REQUEST: username or id_req doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.Error("error getting requesting id or username2")
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

	// Create a variable searcher and assign the value of id and username
	var user components.User
	user.IdUser.Id = id_req
	user.Usname = usname

	// Take following list
	followers, err := rt.db.GetFollowed(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// Check if there are following --> len(following list) != 0
	if len(followers) == 0 {
		// doesn't have followers so write in the header Status No Content and return
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// initialize the variable photos with type as list of PostedPhoto
	var photos []components.PostedPhoto
	// len(following list) != 0 --> iterate the list
	for _, follower := range followers {
		// get following photos
		follPhoto, err := rt.db.GetStream(follower)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(components.InternalServerError))

			if err != nil {
				ctx.Logger.WithError(err).Error("error writing response")
			}

			ctx.Logger.WithError(err).Error("error checking photo existence, details")
			return
		}

		// iterate photo list and append for each photo in the variable photos with type as list of PostedPhoto
		for i := 0; i < len(follPhoto); i++ {
			photos = append(photos, follPhoto[i])
		}
	}

	// if len(photos) == 0 then write in the header status no content and return
	if len(photos) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// status 200
	w.WriteHeader(http.StatusOK)
	// Send the output to the user
	_ = json.NewEncoder(w).Encode(photos)

}
