package routes

import (
	"BE-Sosmed/features/comments"
	"BE-Sosmed/features/postings"
	"BE-Sosmed/features/users"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, uh users.Handler, ph postings.Handler, ch comments.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	routeUser(e, uh)
	routePosting(e, ph)
	routeComment(e, ch)
}

func routeUser(e *echo.Echo, uh users.Handler) {
	e.POST("/users", uh.Register())
	e.POST("/login", uh.Login())
	e.GET("/users/:id", uh.ReadById(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PUT("/users", uh.Update(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/users", uh.Delete(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/users/:username", uh.ReadByUsername())
}

func routePosting(e *echo.Echo, ph postings.Handler) {
	e.POST("/posts", ph.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/posts", ph.GetAll())
	e.GET("/post/:id", ph.GetByPostID())
	e.GET("/posts/:username", ph.GetByUsername())
	e.PUT("/posts/:id", ph.Update(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/posts/:id", ph.Delete(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.POST("/posts/:id", ph.LikePost(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func routeComment(e *echo.Echo, ch comments.Handler) {
	e.POST("/comments", ch.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/comments/:commentId", ch.Delete(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PUT("/comments/:commentId", ch.Update(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}
