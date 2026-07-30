package main

import "github.com/go-chi/chi/v5"

type GreetingOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func main() {
	// router
	router := chi.NewMux()
}
