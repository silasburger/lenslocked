package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/silasburger/lenslocked/context"
	"github.com/silasburger/lenslocked/errors"
	"github.com/silasburger/lenslocked/models"
)

type Galleries struct {
	Templates struct {
		New   Template
		Edit  Template
		Index Template
		Show  Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID    int
		Title     string
		Published bool
	}
	user := context.User(r.Context())
	data.UserID = user.ID
	data.Title = r.FormValue("title")
	data.Published = false
	gallery, err := g.GalleryService.Create(data.Title, data.UserID, data.Published)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID        int
		Title     string
		Published bool
	}
	gallery := context.Gallery(r.Context())
	// gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	// if err != nil {
	// 	return
	// }
	data.ID = gallery.ID
	data.Title = gallery.Title
	data.Published = gallery.Published
	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery := context.Gallery(r.Context())
	// gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	// if err != nil {
	// 	return
	// }
	gallery.Title = r.FormValue("title")
	if r.FormValue("published") == "" {
		gallery.Published = false
	} else if r.FormValue("published") == "true" {
		gallery.Published = true
	}
	err := g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID        int
		Title     string
		Published bool
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{ID: gallery.ID, Title: gallery.Title, Published: gallery.Published})
	}
	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery := context.Gallery(r.Context())
	// if err != nil {
	// 	return
	// }
	var data struct {
		ID     int
		Title  string
		Images []string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	data.Images = createImages()
	g.Templates.Show.Execute(w, r, data)
}

func createImages() []string {
	var images []string
	for i := 0; i < 20; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		catImageURL := fmt.Sprintf("https://placedog.net/%d/%d", w, h)
		images = append(images, catImageURL)
	}
	return images
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery := context.Gallery(r.Context())
	// gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	// if err != nil {
	// 	return
	// }
	err := g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

// func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
// 	id, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusNotFound)
// 	}
// 	gallery, err := g.GalleryService.ByID(id)
// 	if err != nil {
// 		if errors.Is(err, models.ErrNotFound) {
// 			http.Error(w, "Gallery not found", http.StatusInternalServerError)
// 			return nil, err
// 		}
// 		http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 		return nil, err
// 	}
// 	for _, opt := range opts {
// 		err := opt(w, r, gallery)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return gallery, err
// }

// func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
// 	user := context.User(r.Context())
// 	if user.ID != gallery.UserID {
// 		http.Error(w, "You are not authorized to access this gallery", http.StatusForbidden)
// 		return fmt.Errorf("user does not have access to this gallery")
// 	}
// 	return nil
// }

// func mustOwnUnpublishedGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
// 	if !gallery.Published {
// 		user := context.User(r.Context())
// 		if user == nil {
// 			http.Redirect(w, r, "/signin", http.StatusFound)
// 			return fmt.Errorf("unauthorized to access unpublished gallery")
// 		}
// 		return userMustOwnGallery(w, r, gallery)
// 	}
// 	return nil
// }

type GalleryMiddleware struct {
	GalleryService *models.GalleryService
}

func (gmw GalleryMiddleware) GetGallery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		galleryID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(galleryID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		gallery, err := gmw.GalleryService.ByID(id)
		if err != nil {
			fmt.Println(err)
			if errors.Is(err, models.ErrNotFound) {
				http.Error(w, "Gallery not found", http.StatusInternalServerError)
				return
			}
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		ctx = context.WithGallery(ctx, gallery)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (gmw GalleryMiddleware) RequireUserOwnGallery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gallery := context.Gallery(r.Context())
		user := context.User(r.Context())
		if gallery.UserID != user.ID {
			http.Error(w, "You are not authorized to access this gallery", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (gmw GalleryMiddleware) RequireGalleryPublished(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gallery := context.Gallery(r.Context())
		if !gallery.Published {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
