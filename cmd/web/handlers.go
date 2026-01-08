package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"adotkaya.playground/internal/models"
	"adotkaya.playground/internal/validator"
)

// =============================================================================
// Form Types
// =============================================================================

// SnippetCreateForm represents the form data for creating a snippet
type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// userSignupForm represents the form data for user registration
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// userLoginForm represents the form data for user login
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// =============================================================================
// Public Handlers
// =============================================================================

// ping responds with OK for health checks
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// home displays the homepage with a list of the latest snippets
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

// =============================================================================
// Snippet Handlers
// =============================================================================

// snippetView displays a single snippet
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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

// snippetCreate displays the form for creating a new snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = SnippetCreateForm{
		Expires: 365, // Default to 1 year
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

// snippetCreatePost processes the snippet creation form submission
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Decode form data
	var form SnippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form fields
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// If validation failed, re-display the form with errors
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// Insert snippet into database
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add success flash message and redirect
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// =============================================================================
// User Authentication Handlers
// =============================================================================

// userSignup displays the user signup form
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, http.StatusOK, "signup.tmpl", data)
}

// userSignupPost processes the user signup form submission
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Decode form data
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form fields
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 255), "name", "This field cannot be more than 255 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 255), "email", "This field cannot be more than 255 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// If validation failed, re-display the form with errors
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	// Attempt to create the user
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add success flash message and redirect to login
	app.sessionManager.Put(r.Context(), "flash", "Successfully signed up. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// userLogin displays the user login form
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.tmpl", data)
}

// userLoginPost processes the user login form submission
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Decode form data
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form fields
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	// If validation failed, re-display the form with errors
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// Attempt to authenticate the user
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Renew session token to prevent session fixation attacks
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Store user ID in session
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Redirect to snippet create page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// userLogoutPost logs out the user and clears their session
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Renew session token to prevent session fixation attacks
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Remove authenticated user ID from session
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Add success flash message
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
