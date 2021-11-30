package main

import (
	"github.com/bmizerany/pat"
	"github.com/mkruczek/go-chat-ws/internal/handlers"
	"net/http"
)

func routs() http.Handler {

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	return mux
}
