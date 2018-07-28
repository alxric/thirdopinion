package main

import (
	"io"
	"thirdopinion/internal/pkg/handlers"

	"github.com/alecthomas/template"
	"github.com/labstack/echo"
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

	t := &Template{
		templates: template.Must(template.ParseGlob("web/template/*.html")),
	}

	e.Static("/static", "web/static")
	e.Renderer = t
	e.GET("/create", handlers.Create)
	e.GET("/view", handlers.View)
	e.GET("/view/*", handlers.View)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
