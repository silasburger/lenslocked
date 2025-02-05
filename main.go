package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/silasburger/lenslocked/controllers"
	"github.com/silasburger/lenslocked/migrations"
	"github.com/silasburger/lenslocked/models"
	"github.com/silasburger/lenslocked/templates"
	"github.com/silasburger/lenslocked/views"
)

func main() {
	// Set up database connection
	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg)
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Set up service
	userService := models.UserService{
		DB: db,
	}
	tokenManager := models.TokenManager{}
	sessionService := models.SessionService{
		DB:           db,
		TokenManager: tokenManager,
	}

	// Set up middleware
	csrfKey := "le4ZmpwOA80pSiVU8qLWJkjonlEm2MWZ"
	csrfMw := csrf.Protect([]byte(csrfKey))
	// TODO: set to true before deployment
	csrf.Secure(false)

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// Set up controllers
	usersC := controllers.Users{
		UsersService:   &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signin.gohtml"))
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signup.gohtml"))
	usersC.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "users/me.gohtml"))

	// Ser up router and routes

	r := chi.NewRouter()

	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Use(middleware.Logger)

	tpl := views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.Get("/signup", usersC.New)

	r.Post("/signup", usersC.Create)

	r.Get("/signin", usersC.SignIn)

	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "greeting.gohtml"))
	r.Get("/greeting", controllers.StaticHandler(tpl))

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	//
	fmt.Println("Starting server on 3000...")
	http.ListenAndServe(":3000", r)
}
