package main

import (
	"fmt"
	"net/http"

	"bryce-ziemer/github.com/lenslocked/controllers"
	"bryce-ziemer/github.com/lenslocked/models"
	"bryce-ziemer/github.com/lenslocked/templates"
	"bryce-ziemer/github.com/lenslocked/views"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	//tpl = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	//r.Get("/signup", controllers.FAQ(tpl))
	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml",
		"tailwind.gohtml",
	))

	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml",
		"tailwind.gohtml",
	))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	csrfKey := "q8csbqhhteveaq3y1fww0z4201ffqyfa" // Created with https://www.gigacalculator.com/randomizers/random-alphanumeric-generator.php
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false), // TODO fix before deploying (dont have secure local development)
		) // returns a function
	http.ListenAndServe(":3000", csrfMw(r))

}
