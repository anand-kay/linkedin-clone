package core

import (
	"net/http"

	"github.com/anand-kay/linkedin-clone/handlers"
	"github.com/anand-kay/linkedin-clone/middlewares"
	"github.com/anand-kay/linkedin-clone/router"
	"github.com/anand-kay/linkedin-clone/server"
)

func NewRouter(server *server.Server) *router.Router {
	rootRouter := router.NewRouter()
	rootRouter.Use((&middlewares.Server{server}).AttachServerToReq)
	rootRouter.Use(middlewares.SetCORS)
	rootRouter.SetNotfoundHandler(handlers.NotFound)
	rootRouter.SetHandler(`.*`, "OPTIONS", handlers.Options)

	signupLoginRouter := rootRouter.NewSubrouter(``)
	signupLoginRouter.SetHandler(`\/signup`, "POST", handlers.Signup)
	signupLoginRouter.SetHandler(`\/login`, "POST", handlers.Login)

	authorizeRouter := rootRouter.NewSubrouter(``)
	authorizeRouter.Use(middlewares.Authorize)

	logoutRouter := authorizeRouter.NewSubrouter(`\/logout`)
	logoutRouter.SetHandler(``, "POST", handlers.Logout)

	postRouter := authorizeRouter.NewSubrouter(`\/post`)
	postRouter.SetHandler(`\/create`, "POST", handlers.CreatePost)
	postRouter.SetHandler(`\/posts`, "GET", handlers.FetchAllPosts)
	postRouter.SetHandler(`\/[0-9]+`, "GET", handlers.FetchPostByID)

	connectionsRouter := authorizeRouter.NewSubrouter(`\/connections`)
	connectionsRouter.SetHandler(``, "GET", handlers.FetchConnections)
	connectionsRouter.SetHandler(`\/sendreq`, "POST", handlers.SendReq)
	connectionsRouter.SetHandler(`\/acceptreq`, "POST", handlers.AcceptReq)
	connectionsRouter.SetHandler(`\/revokereq`, "DELETE", handlers.RevokeReq)

	userRouter := authorizeRouter.NewSubrouter(`\/user`)
	userRouter.SetHandler(`\/info`, "GET", handlers.UserInfo)

	return rootRouter
}

func GetHandler(router *router.Router) http.Handler {
	return router.GetHandler()
}
