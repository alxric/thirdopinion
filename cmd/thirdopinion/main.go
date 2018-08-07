package main

import (
	"io"
	"os"
	"thirdopinion/internal/pkg/handlers"

	"github.com/alecthomas/template"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog"
)

// Template definition
type Template struct {
	templates *template.Template
}

// Render a template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var (
	version   = "development"
	branch    = ""
	commit    = ""
	buildUser = ""
	goVersion = ""
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger.Info().Str(
		"version", version,
	).Str(
		"branch", branch,
	).Str(
		"commit", commit,
	).Str(
		"go-version", goVersion,
	).Str("build-user", buildUser).Msg("Starting Third Opinion")
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
