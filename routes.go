package celeritas

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (c *Celeritas) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID) //injects request id to context of each request
	mux.Use(middleware.RealIP)    //inject ip address to request
	if c.Debug {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer) //recovers from panics
	mux.Use(c.SessionLoad)
	mux.Use(c.NoSurf)

	return mux
}
