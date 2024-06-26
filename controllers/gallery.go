package controllers

import (
	"bryce-ziemer/github.com/lenslocked/context"
	"bryce-ziemer/github.com/lenslocked/models"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
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

	// we use a seperate type so we can allow for differnt display format vs model
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string // URL escaped filename
	}

	var data struct {
		ID     int
		Title  string
		Images []Image
	}

	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})

	}
	g.Templates.Show.Execute(w, r, data)

}

// this method shows the view that allows editing of a gallery
func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery) // does error handling
	if err != nil {
		return
	}

	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string // URL escaped filename
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
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})

	}

	g.Templates.Edit.Execute(w, r, data)

}

// this method shows the view that allows editing of a gallery
func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery) // does error handling
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

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	//panic("in controller Del")
	gallery, err := g.galleryByID(w, r, userMustOwnGallery) // does error handling
	if err != nil {
		return
	}

	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "Somethign went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r) // Call this because it gets the basename of the file (so will ignore any user added pathing - more secure)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	image, err := g.GalleryService.Image(galleryID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return

		}
		fmt.Println(err)
		http.Error(w, "soemthign went wrong", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, image.Path)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	// look up gallery and that user owns it
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = r.ParseMultipartForm(5 << 20) // 5mb
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	fileheaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileheaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		err = g.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			var fileErr models.FileError
			if errors.As(err, &fileErr) {
				msg := fmt.Sprintf("%v has an invalid content type or extension", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		http.Error(w, "soemthign went wrong", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
}

type galleryOpt func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
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
		err = opt(w, r, gallery)
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
