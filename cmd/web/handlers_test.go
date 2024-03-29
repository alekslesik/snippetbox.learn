package main

import (
	"bytes"

	"net/http"
	"testing"
)

type EmptyHandler http.Handler

// Test Handlers pattern

// ping() GET /ping
func TestPing(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	// contains a couple of mock loggers (which discard anything written to them).
	// app := &application{
	// 	errorLog: log.New(ioutil.Discard, "", 0),
	// 	infoLog:  log.New(ioutil.Discard, "", 0),
	// }
	app := newTestApplication(t, false)

	// We then use the httptest.NewTLSServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() to shutdown the server when
	// the test finishes.
	// ts := httptest.NewTLSServer(app.routes())
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained
	// in the ts.URL field. We can use this along with the ts.Client().Get()
	// method to make a GET /ping request against the test server. This
	// returns a http.Response struct containing the response.
	// rs, err := ts.Client().Get(ts.URL + "/ping")
	// if err != nil {
	//    t.Fatal(err)
	// }

	statusCode, _, body := ts.get(t, "/ping")

	// Initialize a new httptest.ResponseRecorder.
	// rr := httptest.NewRecorder()

	// Initialize a new dummy http.Request.
	// r, err := http.NewRequest("GET", "/", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// Call the ping handler function, passing in the
	// httptest.ResponseRecorder and http.Request.
	// ping(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the
	// http.Response generated by the ping handler.
	// rs := rr.Result()

	// We can then examine the http.Response to check that the status code
	// written by the ping handler was 200.
	if statusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, statusCode)
	}

	// And we can check that the response body written by the ping handler
	// equals "OK".
	// defer rs.Body.Close()
	// body, _ := io.ReadAll(rs.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

// home() GET /
func TestHome(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies
	app := newTestApplication(t, false)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Create a new instance of our application with errors struct which uses the mocked
	// dependencies.
	appERR := newTestApplicationERR(t)
	// Establish a new test server for running end-to-end tests.
	tsERR := newTestServer(t, appERR.routes())

	testCases := []struct {
		desc     string
		urlPath  string
		err      bool
		wantCode int
		wantBody []byte
	}{
		{
			desc: "Valid", urlPath: "/", err: false, wantCode: http.StatusOK, wantBody: []byte(`<th>Title</th>`),
		},
		{
			desc: "Latest() ERR", urlPath: "/", err: true, wantCode: http.StatusInternalServerError, wantBody: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			switch tC.err {
			case false:
				code, _, body := ts.get(t, tC.urlPath)

				if code != tC.wantCode {
					t.Errorf("want %v, get %v", tC.wantCode, code)
				}

				if !bytes.Contains(body, tC.wantBody) {
					t.Errorf("want body to contain %q", tC.wantBody)
				}
			case true:
				code := tsERR.getERR(t, tC.urlPath)
				if code != tC.wantCode {
					t.Errorf("want %v, get %v", tC.wantCode, code)
				}
			}
		})
	}
}

// about() GET /about
func TestAbout(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies
	app := newTestApplication(t, false)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testCases := []struct {
		desc     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{
			desc: "Valid", urlPath: "/about", wantCode: http.StatusOK, wantBody: []byte(`Lorem ipsum dolor sit amet consectetur`),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			code, _, body := ts.get(t, tC.urlPath)

			if code != tC.wantCode {
				t.Errorf("want %v, get %v", tC.wantCode, code)
			}

			if !bytes.Contains(body, tC.wantBody) {
				t.Errorf("want body to contain %q", tC.wantBody)
			}
		})
	}
}

// loginUserForm() GET /user/login
func TestLoginUserForm(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies
	app := newTestApplication(t, false)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testCases := []struct {
		desc     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{
			desc: "Valid", urlPath: "/user/login", wantCode: http.StatusOK, wantBody: []byte(`<input type="submit" value="Login">`),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			code, _, body := ts.get(t, tC.urlPath)

			if code != tC.wantCode {
				t.Errorf("want %v, get %v", tC.wantCode, code)
			}

			if !bytes.Contains(body, tC.wantBody) {
				t.Errorf("want body to contain %q,\n body %q", tC.wantBody, body)
			}
		})
	}
}

//TODO loginUser() POST /user/login

// TODO createSnippetForm() GET /snippet/create
func TestCreateSnippetForm(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies
	app := newTestApplication(t, true)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testCases := []struct {
		desc     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{
			desc: "Valid", urlPath: "/snippet/create", wantCode: http.StatusOK, wantBody: []byte(`<input type="submit" value="Publish snippet">`),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Run(tC.desc, func(t *testing.T) {
				code, _, body := ts.get(t, tC.urlPath)

				if code != tC.wantCode {
					t.Errorf("want %v, get %v", tC.wantCode, code)
				}

				if !bytes.Contains(body, tC.wantBody) {
					t.Errorf("want body to contain %q", tC.wantBody)
				}
			})
		})
	}
}

//TODO createSnippet() POST /snippet/create

// showSnippet() GET /snippet/:id
func TestShowSnippet(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked // dependencies.
	app := newTestApplication(t, false)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our // application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
		{"Server error", "/snippet/100", http.StatusInternalServerError, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _, body := ts.get(t, tt.urlPath)

			if statusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, statusCode)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}

}

//TODO signupUserForm() GET /user/signup

// signupUser() POST /user/signup
func TestSignupUser(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running and end-to-end test.
	app := newTestApplication(t, false)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the
	// response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	// Log the CSRF token value in our test output.
	t.Log(csrfToken)
}

//TODO logoutUser() POST /user/logout
