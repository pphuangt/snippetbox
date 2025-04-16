package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippedbox/internal/models"
	"time"
)

type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	config         config
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

type config struct {
	addr      string
	staticDir string
	dsn       string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{logger: logger}
	flag.StringVar(&app.config.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.config.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&app.config.dsn, "dsa", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	db, err := openDB(app.config.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	app.snippets = &models.SnippetModel{DB: db}
	app.users = &models.UserModel{DB: db}

	TemplateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	app.templateCache = TemplateCache

	app.formDecoder = form.NewDecoder()

	app.sessionManager = scs.New()
	app.sessionManager.Store = mysqlstore.New(db)
	app.sessionManager.Lifetime = 12 * time.Hour
	app.sessionManager.Cookie.Secure = true

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		},
	}

	server := http.Server{
		Addr:         app.config.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server on", slog.Any("addr", app.config.addr))

	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
