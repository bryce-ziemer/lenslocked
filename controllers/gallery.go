package controllers

import (
	"bryce-ziemer/github.com/lenslocked/context"
	"bryce-ziemer/github.com/lenslocked/models"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		Show  Template
		New   Template
		Edit  Template
		Index Template
	}

	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}

	data.Title = r.FormValue("title") // possible to grab title key from url query params...?
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}

	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}

	// If gallery created correctly send use to the edit page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)

}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r) // does error handling
	if err != nil {
		return
	}

	var data struct {
		ID     int
		Title  string
		Images []string
	}

	data.ID = gallery.ID
	data.Title = gallery.Title
	for i := 0; i < 20; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		catImageURL := fmt.Sprintf("https://placekitten.com/%d/%d", w, h)
		data.Images = append(data.Images, catImageURL)
	}
	g.Templates.Show.Execute(w, r, data)

}

// this method shows the view that allows editing of a gallery
func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery) // does error handling
	if err != nil {
		return
	}

	// MOIVED to userMustOwnGallery
	// // check that user is allowed to edit. only those allowed can do so
	// user := context.User(r.Context())
	// if gallery.UserID != user.ID {
	// 	http.Error(w, "You are not authorized to edi thtis gallery", http.StatusForbidden)
	// 	return
	// }

	// show view - view displays the 'editable' gallery
	// this is what we pass to the view. We do not use gallery, because in the general case there maybe
	// info we do not want to pass
	// (this is why we create a new struct)
	var data struct {
		ID    int
		Title string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	g.Templates.Edit.Execute(w, r, data)

}

// this method shows the view that allows editing of a gallery
func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery()) // does error handling
	if err != nil {
		return
	}

	// MOIVED to userMustOwnGallery
	// check that user is allowed to edit. only those allowed can do so
	// user := context.User(r.Context())
	// if gallery.UserID != user.ID {
	// 	http.Error(w, "You are not authorized to edi thtis gallery", http.StatusForbidden)
	// 	return
	// }

	// update the gallery
	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// go back to edit page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)

}

// look up all of a users galelries
// send to a template soi can render them in some way
func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	// make new data struct to be clear about what sending in, and nopt send any data that is not needed
	// not bad to to send full model
	// * seems similar to vm in mvvm *
	type Gallery struct {
		ID    int
		Title string
	}

	var data struct {
		Galleries []Gallery
	}

	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserId(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// image if had a timestamp - we can do thing at this stage to convert to helper text
	// seems like VM
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}
	g.Templates.Index.Execute(w, r, data)
}

type galleryOpt func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...) (*models.Gallery, error) {
	// get id of gallery to work with
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}

	// get gallery to edit
	gallery, err := g.GalleryService.ByID(id)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return nil, err
	}


	for _, opt := range opts {
		err = opt(w,r,gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	// check that user is allowed to edit. only those allowed can do so
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edi thtis gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}

	return nil

}
