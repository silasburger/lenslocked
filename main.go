package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/silasburger/lenslocked/controllers"
	"github.com/silasburger/lenslocked/templates"
	"github.com/silasburger/lenslocked/views"
)

type PGConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (c PGConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

func main() {
	pgconfig := PGConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLMode:  "disabled",
	}

	sql.Open("pgx", pgconfig.String())

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tpl := views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	UsersC := controllers.Users{}
	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signup.gohtml"))
	UsersC.Templates.New = tpl
	r.Get("/signup", UsersC.New)

	r.Post("/signup", UsersC.Create)

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "greeting.gohtml"))
	r.Get("/greeting", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting server on 3000...")
	http.ListenAndServe(":3000", r)
}
