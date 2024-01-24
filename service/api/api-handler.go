package api

import (
	"github.com/julienschmidt/httprouter"
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
	rt.router.GET("/photos", rt.wrap(rt.getUserPhotos))
	rt.router.POST("/photos", rt.wrap(rt.postPhoto))
	rt.router.GET("/photos/:photoID", rt.wrap(rt.getPhoto))
	rt.router.DELETE("/photos/:photoID", rt.wrap(rt.deletePhoto))
	rt.router.PUT("/photos/:photoID/likes/:userID", rt.wrap(rt.putUserLike))
	rt.router.DELETE("/photos/:photoID/likes/:userID", rt.wrap(rt.deleteUserLike))
	rt.router.POST("/photos/:photoID/comments", rt.wrap(rt.postComment))
	rt.router.DELETE("/photos/:photoID/comments/:commentID", rt.wrap(rt.deleteComment))
	rt.router.GET("/photos/:photoID/comments", rt.wrap(rt.getComments))
	rt.router.GET("/users/:userID/stream", rt.wrap(rt.getUserStream))

	fileServer := http.FileServer(http.Dir("./uploads"))
	serveFiles := http.StripPrefix("/uploads", fileServer).ServeHTTP

	rt.router.GET("/uploads/*filepath", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		serveFiles(w, r)
	})

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
