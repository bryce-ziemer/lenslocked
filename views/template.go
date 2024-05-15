package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"bryce-ziemer/github.com/lenslocked/context"
	"bryce-ziemer/github.com/lenslocked/models"

	"github.com/gorilla/csrf"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(path.Base(patterns[0])) // use path.base to pass tpl.ParseFS file naming...?

	// Need to define function before calling ParseFS
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented") // placeholder so hen parse template do not get error
			},

			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented") // placeholder so hen parse template do not get error
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil

}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone() // avoid template race condition (race condition occurs because pointer to template maye chnage when multiple users make requests near in time to one another)
	// We solve the race condition by cloning the template on a per user/request basis

	if err != nil {
		log.Printf("Cloning template: %v", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r) // update place holder in ParseFS (bc we need the request and did not want request to be in ParseFS)
			},

			"currentUser": func() *models.User {
				return context.User(r.Context())
			},

			"errors": func() []string {
				return errMsgs

			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer // in memory now - may lower performance if rendering large templates
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "there was an error executing the template.", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, err.Error())
		}
	}

	return msgs

}
