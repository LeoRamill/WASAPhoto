package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

//   /users/{user-id}/profile/following/{following-id}: PUT and DELETE

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user-id")        // owner user_id
	foll_usn := ps.ByName("following-id")   // id of banned
	id_req := r.Header.Get("Authorization") // authorize user_id

	// check if exist the username and foll_usn
	if username == "" || foll_usn == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting photo id")
		return
	}

	// Bad Request 400 --> length username and foll_usn
	if (len(username) < 6 && len(username) > 12) || (len(foll_usn) < 6 && len(foll_usn) > 12) {
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

	// get foll_id by foll username
	foll_id, err := rt.db.GetUserID(components.Username{Usname: foll_usn})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		return
	}
	// check if exist the foll_id
	if foll_id == "" {
		// STATUS NOT FOUND: id doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	var code_foll components.UserID
	err = json.NewDecoder(r.Body).Decode(&code_foll)
	// Now check if error exist --> error different to nil
	if err != nil {
		// The body was not parsable JSON --> reject it
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error parsing request body, details: %w", err).Error())

		return
	}

	// Check if in the decode body there is the code-user as equal to foll-id
	if code_foll.IdUser.Id != foll_id {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error parsing request body, details: %w", err).Error())

		return
	}

	// CREATE VARIABLE FOLLOWING
	var following components.User
	following.IdUser.Id = foll_id
	following.Usname = foll_usn

	// Check if it appears/exists in the database
	user_check, err := rt.db.CheckUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error validating user")
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

	// insert the follow relationship into the database
	// banned check
	banned_check, err := rt.db.CheckBanned(user, following)
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// if it's true, write header status forbidden because I can't follow people that were banned by me
	if banned_check {
		// User was banned by owner
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(components.ForbiddenError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// insert follow in the database
	_, err = rt.db.FollowUser(user, following)

	if err != nil {
		ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
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

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("user-id")        // owner user_id
	foll_usn := ps.ByName("following-id")   // id of banned
	id_req := r.Header.Get("Authorization") // authorize user_id

	// check if exist the username and foll_usn
	if username == "" || foll_usn == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.Error("error getting photo id")
		return
	}

	// Bad Request 400 --> length username and foll_usn
	if (len(username) < 6 && len(username) > 12) || (len(foll_usn) < 6 && len(foll_usn) > 12) {
		// BAD REQUEST: username or banned doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// get userID by photoauthorusername
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

	// get foll_id by foll username
	foll_id, err := rt.db.GetUserID(components.Username{Usname: foll_usn})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error getting user")
		return
	}
	// check if exist the foll_id
	if foll_id == "" {
		// STATUS NOT FOUND: id doesn't exist in the db
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(components.NotFoundError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	var following components.User
	following.IdUser.Id = foll_id
	following.Usname = foll_usn

	// Check if it appears/exists in the database
	user_check, err := rt.db.CheckUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error validating user")
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

		ctx.Logger.WithError(err).Error("invalid token")
		return
	}

	// insert the follow relationship into the database
	// banned check
	banned_check, err := rt.db.CheckBanned(user, following) // Fare Check banner
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}
	// if it's true, write header status forbidden because I can't follow people that were banned by me
	if banned_check {
		// User was banned by owner
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte(components.ForbiddenError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// delete follow in the database
	_, err = rt.db.UnFollowUser(user, following)

	if err != nil {
		ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
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
