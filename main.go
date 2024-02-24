package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1?")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:bryce.ziemer@gmail.com\">bryce.ziemer@gmail.com</a>.")
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

	default:
		//w.WriteHeader(http.StatusNotFound)
		//fmt.Fprint(w, "Page Not Found!")
		http.Error(w, "Page Not Found", http.StatusNotFound)

	}

}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	//http.ListenAndServe(":3000", http.HandlerFunc(pathHandler))
	http.ListenAndServe(":3000", router)

}
