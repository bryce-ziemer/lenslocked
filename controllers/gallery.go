package controllers

import (
	"bryce-ziemer/github.com/lenslocked/models"
	"net/http"
)

type Galleries struct {
	Templates struct {
		New Template
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
