package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"adotkaya.playground/internal/models"
	"adotkaya.playground/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type SnippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Use the new render helper.
	data := app.newTemplateData(r)
	data.Snippets = snippets
	// Pass the data to the render() helper as normal.
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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
	// Use the new render helper.
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	{
		data := app.newTemplateData(r)

		data.Form = SnippetCreateForm{
			Expires: 365,
		}

		app.render(w, http.StatusOK, "create.tmpl", data)
	}
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// ParseForm() is limited by 10mb, if need more; use http.MaxBytesReader() before ParseForm().
	// ParseForm() wont raise a flag if exceed. So there will be no logs and just bad user experience...
	// MaxBytesReader sets flag on http.ResponseWriter which instructs server to close the TCP connection
	// ex: r.Body = http.MaxBytesReader(w,r.Body,4096) then keep going (4096 bytes example code).
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := SnippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// Validate if not empty
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank.")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must either 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
