package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Session routes
	// /session
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Search routes
	// /users
	rt.router.GET("/users/:user-id/", rt.wrap(rt.searchUser))

	// HOMEPAGE: Stream routes
	// /users/{user-id}/homepage:
	rt.router.GET("/users/:user-id/homepage", rt.wrap(rt.getStream))

	// /users/{user-id}/homepage/{photo-id}/likes/{like-id}:
	rt.router.PUT("/users/:user-id/homepage/:photo-id/likes/:like-id", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:user-id/homepage/:photo-id/likes/:like-id", rt.wrap(rt.unlikePhoto))

	// MODIFICA IMPORTANTE
	// /users/{user-id}/homepage/{photo-id}/comments
	rt.router.POST("/users/:user-id/homepage/:photo-id/comments", rt.wrap(rt.createComment))
	rt.router.DELETE("/users/:user-id/homepage/:photo-id/comments/:comment-id", rt.wrap(rt.deleteComment))

	// PROFILE: info profile
	// /users/{user-id}/profile PUT
	// Update Username
	rt.router.PUT("/users/:user-id/profile", rt.wrap(rt.updateUsername))

	rt.router.GET("/users/:user-id/profile", rt.wrap(rt.getUserProfile))

	// PROFILE: banned
	// /users/{user-id}/profile/banned:
	rt.router.GET("/users/:user-id/profile/banned", rt.wrap(rt.getBans))
	// /users/{user-id}/profile/banned/{banned-id}:
	rt.router.PUT("/users/:user-id/profile/banned/:banned-id", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:user-id/profile/banned/:banned-id", rt.wrap(rt.unBanUser))

	// PROFILE: following
	// /users/{user-id}/profile/following/{following-id}:
	rt.router.PUT("/users/:user-id/profile/following/:following-id", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:user-id/profile/following/:following-id", rt.wrap(rt.unfollowUser))

	// PROFILE: photos
	// /users/{user-id}/profile/photos:
	rt.router.POST("/users/:user-id/profile/photos", rt.wrap(rt.postPhoto))

	// /users/{user-id}/profile/photos/{photo-id}:
	// rt.router.GET("/users/:user-id/profile/photos/:photo-id", rt.wrap(rt.getOwnerPostedPhoto))
	rt.router.GET("/users/:user-id/profile/photos/:photo-id", rt.wrap(rt.getPhoto))
	rt.router.DELETE("/users/:user-id/profile/photos/:photo-id", rt.wrap(rt.deletePostedPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
