package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alekslesik/snippetbox.learn/pkg/models"
	"github.com/alekslesik/snippetbox.learn/pkg/models/mysql"
	"github.com/golangcollege/sessions"

	_ "github.com/go-sql-driver/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	gopath   string
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	snippets interface {
		Insert(title, content, expires string) (int, error)
		Get(id int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	users         interface {
		Insert(name, email, password string) error
		Authenticate(email, password string) (int, error)
		Get(id int) (*models.User, error)
	}
}

func main() {
	// Command-line flag parsing
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	dsn := flag.String("dsn", "web:ndJMv9zrJw@/snippetbox?parseTime=true", "Название MySQL источника данных")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")
	flag.Parse()

	// Go path
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		log.Fatal("GOPATH is not defined")
	}

	// Open/create log file
	f, err := os.OpenFile(gopath + "/src/github.com/alekslesik/snippetbox.learn/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Loggers
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open DB connection pull
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialise new cache pattern
	templateCache, err := newTemplateCache(gopath + "/src/github.com/alekslesik/snippetbox.learn/ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new session manager
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	// Initialisation application struct
	app := &application{
		gopath:        gopath,
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings the server to use
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:               tls.VersionTLS12,
	}

	// Confugure server
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Server started on http://127.0.0.1%s", *addr)
	fmt.Printf("Server started on http://127.0.0.1%s", *addr)

	// Use the ListenAndServeTLS() method to start the HTTPS server. We
	// pass in the paths to the TLS certificate and corresponding private key a
	// the two parameters.
	err = srv.ListenAndServeTLS(gopath+"/src/github.com/alekslesik/snippetbox.learn/tls/cert.pem", gopath+"/src/github.com/alekslesik/snippetbox.learn/tls/key.pem")

	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
