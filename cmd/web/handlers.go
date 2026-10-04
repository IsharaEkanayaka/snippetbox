package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox/internal/models"
)

func(app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//Initialize a slice containing the paths to the 2 files. It's important
	// to note that the file containing our base template must be the "first"
	//file in the slice.
	files := []string{
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/base.tmpl",
	}

	//use the template.ParseFiles() function to read the files and store the 
	//template in a template set.Notice that we can pass the slice of file
	//paths as a variadic parameter?
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base"
	//template as the response body
	err = ts.ExecuteTemplate(w,"base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//create some variables holding dummy data. we'll remove these later on
	title:= "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	//call the SnippetModel.Insert() method to insert a new record in the database
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//redirect the user to the relevant snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}