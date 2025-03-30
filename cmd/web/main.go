package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"net/http"
	"os"
	"snippedbox/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
	config   config
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

	logger.Info("starting server on", slog.Any("addr", app.config.addr))
	err = http.ListenAndServe(app.config.addr, app.routes())
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
