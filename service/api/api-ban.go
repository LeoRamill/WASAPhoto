package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

func (rt *_router) getBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("user-id")
	id_req := r.Header.Get("Authorization")

	// check if exist path name
	if username == "" {
		// Doesn't exist the username or id_req so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		return
	}

	// Bad Request 400 -->  length username or id
	if (len(username) < 6 && len(username) > 12) || (len(id_req) != 64) {
		// BAD REQUEST: username or id_req doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// get userID by usname (author of the photo)
	id_usname, err := rt.db.GetUserID(components.Username{Usname: username})
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

	// Create a variable searcher and assign the value of id and username
	var user components.User
	user.IdUser.Id = id_usname
	user.Usname = username

	// Take bans list
	bans, err := rt.db.GetBans(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	// Send the output to the user.
	_ = json.NewEncoder(w).Encode(bans)
}

// /users/{user-id}/profile/banned/{banned-id}: PUT and DELETE
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user-id")     // owner user_id
	banned_usn := ps.ByName("banned-id") // id of banned

	// check if exist the banned or username
	if username == "" || banned_usn == "" {
		// Doesn't exist the username or banned so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		return
	}

	// Bad Request 400 --> length username
	if (len(username) < 6 && len(username) > 12) || (len(banned_usn) < 6 && len(banned_usn) > 12) {
		// BAD REQUEST: username or banned doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// get userID by usname
	id_usname, err := rt.db.GetUserID(components.Username{Usname: username})
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

	// CREATE VARIABLE
	var user components.User
	user.IdUser.Id = id_usname
	user.Usname = username

	// get  bannedID by banned_usname
	banned_id, err := rt.db.GetUserID(components.Username{Usname: banned_usn})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error validating user")
		return
	}
	// check if exist the banned_id
	if banned_id == "" {
		// STATUS NOT FOUND: id doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error getting banned id")
		return
	}

	var ban components.User
	ban.IdUser.Id = banned_id
	ban.Usname = banned_usn

	// Check if it appears/exists in the database
	user_check, err := rt.db.CheckUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		return
	}
	// if it's false (so doesn't exist in db) then return Status Not Found
	if !user_check {
		// STATUS NOT FOUND: user doesn't exist in the db
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error authenticating")
		return
	}

	// banned check --> se è bannato esci perché non lo posso bannare due volte
	banned_check, err := rt.db.CheckBanned(ban, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		return
	}
	// if it's true, write header status forbidden because I can't ban two times
	if banned_check {
		// User was banned by owner
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(components.ForbiddenError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// insert ban in the database
	_, err = rt.db.BanUser(user, ban)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
		return
	}

	// Se fa parte dei followed --> Unfollow
	follower_check, err := rt.db.CheckFollow(user, ban)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		return
	}
	if follower_check {
		_, err = rt.db.UnFollowUser(user, ban)

		if err != nil {
			ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(components.InternalServerError))

			if err != nil {
				ctx.Logger.WithError(err).Error("error writing response")
			}
			return
		}
	}

	// Se fa parte dei followed --> Unfollow
	following_check, err := rt.db.CheckFollow(ban, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		return
	}
	if following_check {
		_, err = rt.db.UnFollowUser(ban, user)

		if err != nil {
			ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(components.InternalServerError))

			if err != nil {
				ctx.Logger.WithError(err).Error("error writing response")
			}
			return

		}
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}

// /users/{user-id}/profile/banned/{banned-id}:
func (rt *_router) unBanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("user-id")        // owner user_id
	banned_usn := ps.ByName("banned-id")    // id of following
	id_req := r.Header.Get("Authorization") // authorize user_id

	// check if exist the username and banned_usn
	if username == "" || banned_usn == "" {
		// Doesn't exist the username or id_req so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		return
	}

	// Bad Request 400 --> length username and banned_usn
	if (len(username) < 6 && len(username) > 12) || (len(banned_usn) < 6 && len(banned_usn) > 12) {
		// BAD REQUEST: username or banned doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// get userID by username
	id_usname, err := rt.db.GetUserID(components.Username{Usname: username})
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

	// CREATE VARIABLE
	var user components.User
	user.IdUser.Id = id_usname
	user.Usname = username

	// get banned_id by banned username
	banned_id, err := rt.db.GetUserID(components.Username{Usname: banned_usn})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error getting banned id")
		return
	}

	// check if exist the banned_id
	if banned_id == "" {
		// STATUS NOT FOUND: id doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	var ban components.User
	ban.IdUser.Id = banned_id
	ban.Usname = banned_usn

	// Check if it appears/exists in the database
	user_check, err := rt.db.CheckUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		return
	}
	// if it's false (so doesn't exist in db) then return Status Not Found
	if !user_check {
		// STATUS NOT FOUND: user doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// delete ban in the database
	_, err = rt.db.UnBanUser(user, ban)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
		return
	}
	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)

}
