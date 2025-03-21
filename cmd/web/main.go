package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	config config
}

type config struct {
	addr      string
	staticDir string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{logger: logger}
	flag.StringVar(&app.config.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.config.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	logger.Info("starting server on", slog.Any("addr", app.config.addr))
	err := http.ListenAndServe(app.config.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
