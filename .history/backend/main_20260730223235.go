package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humachi"
	"github.com/go-chi/chi/v5"
)

type GreetingOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func main() {
	// router
	router := chi.NewMux()

	api := humachi.New(router, huma.DefaultConfig("ICTSC Score Server"))

	// get
	huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"John" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	http.ListenAndServe(":8080", router)

}
