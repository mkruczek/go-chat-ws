package main

import "net/http"

func main() {

	mux := routs()

	http.ListenAndServe(":8888", mux)
}
