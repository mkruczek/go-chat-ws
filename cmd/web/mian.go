package main

import (
	"github.com/mkruczek/go-chat-ws/internal/handlers"
	"net/http"
)

func main() {

	mux := routs()

	go handlers.ListenForWsChannel()

	http.ListenAndServe(":8888", mux)
}
