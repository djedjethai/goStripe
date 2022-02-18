package main

// START THE APP WITH: air (in the root dir)

// go path directory /home/jerome/Documents/code/go/domainhex/ but no package there ..??

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/djedjethai/goStripe/internal/driver"
	"github.com/djedjethai/goStripe/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version    = "1.0.0"
	cssVersion = "1"
)

var session *scs.SessionManager

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	secretKey string
	frontend  string
}

type application struct {
	config      config
	infoLog     *log.Logger
	errorLog    *log.Logger
	templeCache map[string]*template.Template
	version     string
	DB          models.DBModel
	Session     *scs.SessionManager
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	app.infoLog.Println(fmt.Sprintf("Starting http server on mode %s on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()
}

func main() {
	// we register a type with the gob encoding
	// we need to add a second {} ()
	// with that we can put the type map[string]interface{} into the session
	// ex (before passing the TransactionData type): gob.Register(map[string]interface{}{})
	gob.Register(TransactionData{})
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment{development|production}")
	// read the connection string(dsn) from the flags(i can change to var env if i prefer)
	flag.StringVar(&cfg.db.dsn, "dsn", "mariadb:password@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "DSN")
	// this is the backend api
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")
	flag.StringVar(&cfg.secretKey, "secret", "khgfhjgfjhgfytfcb", "secret key")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to frontend")

	flag.Parse()
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)

	// setup the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = mysqlstore.New(conn)

	if err != nil {
		// use .Fatal() to exit as the db connection failed
		errorLog.Fatal(err)
	}
	defer conn.Close()

	tc := make(map[string]*template.Template)

	app := &application{
		config:      cfg,
		infoLog:     infoLog,
		errorLog:    errorLog,
		templeCache: tc,
		version:     version,
		DB:          models.DBModel{DB: conn},
		Session:     session,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
