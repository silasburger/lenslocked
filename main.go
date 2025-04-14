package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/silasburger/lenslocked/controllers"
	"github.com/silasburger/lenslocked/migrations"
	"github.com/silasburger/lenslocked/models"
	"github.com/silasburger/lenslocked/templates"
	"github.com/silasburger/lenslocked/views"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL, err = models.DefaultPostgresConfig()
	if err != nil {
		return cfg, err
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	secureStr := os.Getenv("CSRF_SECURE")
	secure, err := strconv.ParseBool(secureStr)
	if err != nil {
		return cfg, err
	}
	cfg.CSRF.Secure = secure

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	// Set up database connection
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Set up service
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB:           db,
		TokenManager: &models.TokenManager{},
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleriesService := &models.GalleryService{
		DB: db,
	}

	// Set up middleware
	// TODO: set to true before deployment
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	// Set up controllers
	usersC := controllers.Users{
		UsersService:         userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signin.gohtml"))
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signup.gohtml"))
	usersC.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "users/me.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "forgot-pw.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "check-your-email.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "reset-pw.gohtml"))
	usersC.Templates.SendSignin = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "send-signin.gohtml"))
	usersC.Templates.EditEmail = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "edit-email.gohtml"))

	galleriesC := controllers.Galleries{
		GalleryService: galleriesService,
	}
	galleriesC.Templates.New = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "galleries/new.gohtml"))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "galleries/edit.gohtml"))
	galleriesC.Templates.Index = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "galleries/index.gohtml"))
	galleriesC.Templates.Show = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "galleries/show.gohtml"))

	// Set up router and routes
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

	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	// TODO: change name to passwordless-signin
	r.Get("/send-signin", usersC.SendSignin)
	r.Post("/send-signin", usersC.ProcessSendSignin)

	// TODO: change name to passwordless-signin-link
	r.Get("/email-signin", usersC.ProcessEmailSignin)

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "greeting.gohtml"))
	r.Get("/greeting", controllers.StaticHandler(tpl))

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.Route("/users/edit-email", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.EditEmail)
		r.Post("/", usersC.ProcessEditEmail)
	})

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting server on 3000...")
	http.ListenAndServe(cfg.Server.Address, r)
}
