package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.postSession))
	rt.router.GET("/users", rt.wrap(rt.getUsers))
	rt.router.GET("/users/:userID", rt.wrap(rt.getUser))
	rt.router.PUT("/users/:userID/username", rt.wrap(rt.putUserUsername))
	rt.router.PUT("/users/:userID/banned/:bannedUserID", rt.wrap(rt.putUserBan))
	rt.router.DELETE("/users/:userID/banned/:bannedUserID", rt.wrap(rt.deleteUserBan))
	rt.router.PUT("/users/:userID/following/:followedID", rt.wrap(rt.putUserFollow))
	rt.router.DELETE("/users/:userID/following/:followedID", rt.wrap(rt.deleteUserFollow))
	rt.router.GET("/users/:userID/followers", rt.wrap(rt.getUserFollowers))
	rt.router.GET("/users/:userID/following", rt.wrap(rt.getUserFollowings))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
