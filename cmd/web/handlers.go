package main

import (
	"errors"
	"fmt"

	"github.com/alekslesik/snippetbox.learn/pkg/forms"
	"github.com/alekslesik/snippetbox.learn/pkg/models"

	"net/http"
	"strconv"
)

// Ping GET /ping
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// Home page GET /
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{
		Snippets: s,
	})
}

// About page GET /about
func (app *application) about(w http.ResponseWriter, r *http.Request)  {
	app.render(w, r, "about.page.html", &templateData{})
}

// Login user GET /user/login
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &templateData{
		Form: forms.New(nil),
	})
}

// Login user POST /user/login
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.html", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Add the ID of the current user to the session
	app.session.Put(r, "userID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// Create snippet GET /snippet/create
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", &templateData{
		// pass a new empty forms.Form object to the template
		Form: forms.New(nil),
	})
}

// Create snippet POST /snippet/create
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// adds any data in POST request bodies to the r.PostForm map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create forms.Form containing the POSTed data from the form
	form := forms.New(r.PostForm)
	// Use validation functions
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// if any errors, redisplay the create.page.html paasingvalidation errors and
	// previously submitted r.PostForm data
	if !form.Valid() {
		app.render(w, r, "create.page.html", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add a string value and key to the session data
	app.session.Put(r, "flash", "Snippet sucessfully created")

	// GET
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

// Show snippet GET /snippet/:id
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.html", &templateData{
		Snippet: s,
	})
}

// Sign up user GET /user/signup
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &templateData{
		Form: forms.New(nil),
	})
}

// Sign up user POST /user/signup
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.html", &templateData{
			Form: form,
		})
		return
	}

	// Try to create a new user record in the database. If the email already exist
	// add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.html", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Otherwise add a confirmation flash message to the session confirming
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	// GET
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Logout user POST /user/logout
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {

	// Remove iserID from session.
	app.session.Remove(r, "userID")
	// Add flash to session.
	app.session.Put(r, "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}