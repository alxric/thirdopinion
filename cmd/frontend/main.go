package main

import (
	"io"
	"thirdopinion/internal/pkg/handlers"

	"github.com/alecthomas/template"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

// Template definition
type Template struct {
	templates *template.Template
}

// Render a template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	t := &Template{
		templates: template.Must(template.ParseGlob("web/template/*.html")),
	}

	e.Static("/static", "web/static")
	e.Renderer = t
	e.GET("/", handlers.View)
	e.GET("/create", handlers.Create)

	e.GET("/view", handlers.View)
	e.GET("/view/*", handlers.View)

	e.GET("/register", handlers.Register)
	e.GET("/login", handlers.Login)
	e.POST("/login", handlers.NewLogin)
	e.POST("/logout", handlers.Logout)

	// API
	e.POST("/api/session/validate", handlers.ValidateSession)
	e.GET("/api/arguments", handlers.ListArguments)
	e.POST("/api/vote", handlers.Vote)
	e.POST("/api/register", handlers.RegisterUser)
	e.POST("/api/create", handlers.CreateArgument)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
