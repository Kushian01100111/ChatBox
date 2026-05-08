package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/kushian01100111/ChatBox/internal/app/chat"
	"github.com/kushian01100111/ChatBox/internal/config"
	Server "github.com/kushian01100111/ChatBox/internal/http"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	addr := flag.String("addr", ":"+config.Port, "HTTP network address")

	hub := chat.NewHub()
	go hub.Run()

	r := Server.NewHandler(hub)
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
