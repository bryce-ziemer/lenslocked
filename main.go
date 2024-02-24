package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1?")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:bryce.ziemer@gmail.com\">bryce.ziemer@gmail.com</a>.")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ</h1>
	<ul>
		<li><b>Is there a free version?</b> Yes! we offer a free trial for 30 days on any paid plans.</li>
		<li><b>Question 2?</b> Answer 2</li>
		<li><b>How do I contact support?</b> Email us - <a href="mailto:support@email.com">support@lenslocked.com</a></li>
	</ul>`)
}

// func pathHandler(w http.ResponseWriter, r *http.Request) {

// 	switch r.URL.Path {
// 	case "/":
// 		homeHandler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)

// 	default:
// 		//w.WriteHeader(http.StatusNotFound)
// 		//fmt.Fprint(w, "Page Not Found!")
// 		http.Error(w, "Page Not Found", http.StatusNotFound)

// 	}

// }

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)

	default:
		//w.WriteHeader(http.StatusNotFound)
		//fmt.Fprint(w, "Page Not Found!")
		http.Error(w, "Page Not Found", http.StatusNotFound)

	}

}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)

}
