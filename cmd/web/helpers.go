package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/alekslesik/snippetbox.learn/pkg/models"
	"github.com/justinas/nosurf"
)

// The serverError helper writes an error message and stack trace to the errorLo
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding descri
// to the user. We'll use this later in the book to send responses like 400 "Bad Request"
// when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// extract pattern depending "name"
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("pattern %s not exist", name))
		return
	}

	// initialize a new buffer
	buf := new(bytes.Buffer)

	// write template to the buffer, instead straight to http.ResponseWriter
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write buffer to http.ResponseWriter
	buf.WriteTo(w)
}

// Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we're not using the *http.Request parameter at the
// moment, but we will do later in the book.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	// Add current time.
	td.CurrentYear = time.Now().Year()
	// Add flash message.
	td.Flash = app.session.PopString(r, "flash")
	// Check if user is authenticate.
	td.AuthenticatedUser = app.authenticatedUser(r)
	// Add the CSRF token to the templateData struct.
	td.CSRFToken = nosurf.Token(r)
	//

	return td
}

// Return userID ID from session
func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
