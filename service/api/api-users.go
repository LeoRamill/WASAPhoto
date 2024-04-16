package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"wasaphoto.uniroma1.it/photo1984766/service/api/reqcontext"
	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// /users/{user-id}/
func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get UserID by the path
	username := ps.ByName("user-id")
	// Get the id by Authorization field
	id_req := r.Header.Get("Authorization")
	// check if exist the id
	if id_req == "" {
		// NO AUTHENTICATION: Doesn't exist the id, so the action to search the username is unauthorized
		// write in the header that the status is unauthorized and return it
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// check if exist path name
	if username == "" {
		// Doesn't exist the username so write in the header Bad Request Status and return it
		// of course username doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		return
	}

	// Bad Request 400 --> length username
	if (len(username) < 6 && len(username) > 12) || (len(id_req) != 64) {
		// BAD REQUEST: username or id_req doesn't satisfy the costraints setting to the api documentation
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.WithError(err).Error(fmt.Sprintf("error checking photo existence, details 2: %s", id_req))

		return
	}

	// get userID by username in the path
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
		// NO AUTHENTICATION: Doesn't exist the id, so the action to search the username is unauthorized
		// write in the header that the status is unauthorized and return it
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(components.UnauthorizedError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		ctx.Logger.Error("error authenticating")
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
	var searcher components.User
	searcher.IdUser.Id = id_req
	searcher.Usname = username

	// Get the username to search
	UsernameToSearch := r.URL.Query().Get("searched-id")
	// check if exist the username to search
	if UsernameToSearch == "" {
		// USERNAME INVALID: write in the header that the status is a bad request since it's an empty string
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error(fmt.Sprintf("error checking photo existence, details 3: %s", err))
		return
	}

	// initialize and assign the variable user with list User type using the SearchUser function
	users, err := rt.db.SearchUser(UsernameToSearch)
	// repsonses: 500
	if err != nil {
		// In this case, there's an error coming from the database. Return an empty json
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("Database has encountered an error")
		_ = json.NewEncoder(w).Encode([]components.User{})
		return

	}

	// check if we found the users
	if len(users) == 0 {
		// if users list is empty, write in header Status No Content and return
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// if the users list is not empty --> iteration
	for _, userToSearch := range users {
		// create a boolean variable to check if the searcher is banned by the searched
		var banned_check1 bool
		banned_check1, err = rt.db.CheckBanned(searcher, userToSearch) // Fare Check banner
		if err != nil {
			ctx.Logger.WithError(err).Error("db.CheckBanned: error executing query")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(components.InternalServerError))

			if err != nil {
				ctx.Logger.WithError(err).Error("error writing response")
			}
			return
		}
		// check if it's true
		if banned_check1 {
			// FORBIDDEN ACTION: if it's true, write in the header that the status is forbidden and then return
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte(components.ForbiddenError))

			if err != nil {
				ctx.Logger.WithError(err).Error("error writing response")
			}
			return
		}

		// create a boolean variable to check if the searched is banned by the searcher
		var banned_check2 bool
		banned_check2, err = rt.db.CheckBanned(userToSearch, searcher) // Fare Check banner
		if err != nil {
			ctx.Logger.WithError(err).Error("db.CheckBanned: error executing query")
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(components.InternalServerError))

			if err != nil {
				ctx.Logger.WithError(err).Error("error writing response")
			}
			return
		}
		// check if it's true
		if banned_check2 {
			// if it's true write status no content and return it
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// if there are user in the user that are not banned write status OK and return the users
	w.WriteHeader(http.StatusOK)
	// Send the output to the user. Instead of giving null for no matches return and empty slice of photos. ( ontrollaerrore)
	_ = json.NewEncoder(w).Encode(users)

}

// Update Username
// /users/{user-id}/profile
func (rt *_router) updateUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// get username
	usname := ps.ByName("user-id")
	/* AUTHORIZATION
	id will contain the value of the "Authorization" header, which often holds some form of authentication token or credentials.
	get requesting (UserID)
	*/
	id_req := r.Header.Get("Authorization")
	// check if exist the id or username
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
		return
	}

	// check if exist the id
	if id_usname == "" {
		w.WriteHeader(http.StatusNotFound)
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

	// get the new username from the request body and read the request body
	var new_usname components.Username
	err = json.NewDecoder(r.Body).Decode(&new_usname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(components.BadRequestError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}

		ctx.Logger.WithError(err).Error("error decoding request body")
		return
	}

	// Create a variable searcher and assign the value of id and username
	var user components.User
	user.IdUser.Id = id_req
	user.Usname = usname

	// Change folder path (user's directories locally)
	oldPath := filepath.Join(photoFolder, usname)
	newPath := filepath.Join(photoFolder, new_usname.Usname)

	// Rename the folder
	err = os.Rename(oldPath, newPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(components.InternalServerError))

		if err != nil {
			ctx.Logger.WithError(err).Error("error writing response")
		}
		return
	}

	// change the username in the database
	ret_username, err := rt.db.UpdateUsername(user, new_usname.Usname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Sprintf("error checking photo existence, details 2: %s", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	// Send the output to the user.
	_ = json.NewEncoder(w).Encode(ret_username)

}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	username := ps.ByName("user-id") // owner user_id

	// check if exist the username
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

	// Bad Request 400 --> length username
	if len(username) < 6 && len(username) > 12 {
		// BAD REQUEST: username doesn't satisfy the costraints setting to the api documentation
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

	// CREATE VARIABLE
	var user components.User
	user.IdUser.Id = id_usname
	user.Usname = username

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

	// Define a List
	// vedere anche questo qualora il problema persiste
	// var followers []components.User
	// var following []components.User
	// var photos []components.PostedPhoto

	// get list of followers
	followers, err := rt.db.GetFollowers(user)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get list of following
	following, err := rt.db.GetFollowed(user)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetFollowing: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get list of photos
	photos, err := rt.db.GetGallery(user)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetPhotosList: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write in the header status OK
	w.WriteHeader(http.StatusOK)
	// Send the output to the user
	_ = json.NewEncoder(w).Encode(components.CompleteProfile{
		Name:      id_usname,
		Nickname:  username,
		Followers: followers,
		Following: following,
		Posts:     photos,
	})

}
