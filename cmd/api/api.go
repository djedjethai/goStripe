package main

import (
	"flag"
	"fmt"
	"github.com/djedjethai/goStripe/internal/driver"
	"github.com/djedjethai/goStripe/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version = "1.0.0"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment{development|production|maintenance}")
	// read the connection string(dsn) from the flags(i can change to var env if i prefer)
	flag.StringVar(&cfg.db.dsn, "dsn", "mariadb:password@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "DSN")
	flag.StringVar(&cfg.smtp.host, "smtphost", "localhost", "smtp host")
	flag.StringVar(&cfg.smtp.username, "smtpuser", "", "smtp user")
	flag.StringVar(&cfg.smtp.password, "smtppass", "", "smtp password")
	flag.IntVar(&cfg.smtp.port, "smtpport", 1025, "smtp port")

	flag.Parse()
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		// use .Fatal() to exit as the db connection failed
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       models.DBModel{DB: conn},
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

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
	app.infoLog.Println(fmt.Sprintf("Starting Back end server on mode %s on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()
}
