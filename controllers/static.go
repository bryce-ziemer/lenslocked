package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! we offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "Question 2?",
			Answer:   "Answer 2",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@email.com">support@lenslocked.com</a>`,
		},
		{
			Question: "Where is your office located?",
			Answer:   "Our entire team is remote!",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
