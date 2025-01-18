package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/silasburger/lenslocked/controllers"
	"github.com/silasburger/lenslocked/models"
	"github.com/silasburger/lenslocked/templates"
	"github.com/silasburger/lenslocked/views"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}
	UsersC := controllers.Users{
		UsersService: &userService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tpl := views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signup.gohtml"))
	UsersC.Templates.New = tpl
	r.Get("/signup", UsersC.New)

	r.Post("/signup", UsersC.Create)

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signin.gohtml"))
	UsersC.Templates.New = tpl
	r.Get("/signin", UsersC.SignIn)

	r.Post("/signin", UsersC.ProcessSignIn)

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "greeting.gohtml"))
	r.Get("/greeting", controllers.StaticHandler(tpl))

	r.Get("/users/me", UsersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	csrfKey := "le4ZmpwOA80pSiVU8qLWJkjonlEm2MWZ"
	csrfMw := csrf.Protect([]byte(csrfKey))
	// TODO: set to true before deployment
	csrf.Secure(false)

	fmt.Println("Starting server on 3000...")
	http.ListenAndServe(":3000", csrfMw(r))
}
