package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/silasburger/lenslocked/controllers"
	"github.com/silasburger/lenslocked/templates"
	"github.com/silasburger/lenslocked/views"
)

// func executeTemplate(w http.ResponseWriter, filePath string) {
// 	w.Header().Set("Content-Type", "text/html")
// 	tpl, err := views.Parse(filePath)
// 	if err != nil {
// 		log.Printf("parsing template: %v", err)
// 		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
// 		return
// 	}
// 	tpl.Execute(w, nil)
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	executeTemplate(w, "templates/home.gohtml")
// }

// func contactHandler(w http.ResponseWriter, r *http.Request) {
// 	executeTemplate(w, "templates/contact.gohtml")
// }
// func faqHandler(w http.ResponseWriter, r *http.Request) {
// 	executeTemplate(w, "templates/faq.gohtml")
// }

// func galleriesHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/plain")

// 	fmt.Fprint(w, chi.URLParam(r, "id"))
// }

// // func pathHandler(w http.ResponseWriter, r *http.Request) {
// // 	switch r.URL.Path {
// // 	case "/":
// // 		homeHandler(w, r)
// // 	case "/contact":
// // 		contactHandler(w, r)
// // 	default:
// // 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// // 	}

// // }

// // type Router struct {
// // }

// // func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// // 	switch r.URL.Path {
// // 	case "/":
// // 		homeHandler(w, r)
// // 	case "/contact":
// // 		contactHandler(w, r)
// // 	default:
// // 		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// // 	}
// // }

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tpl := views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.Get("/signup", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "signup.gohtml"))))

	tpl = views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "greeting.gohtml"))
	r.Get("/greeting", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting server on 3000...")
	http.ListenAndServe(":3000", r)
}
