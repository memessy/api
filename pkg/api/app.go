package api

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"memessy-api/pkg/api/meme"
	"net/http"
)

type Config struct {
	StaticDir    http.Dir
	StaticPrefix string
}

func NewApi(
	config Config,
	memeResource meme.Resource,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/memes/", memeResource.Create)
	r.Get("/memes/", memeResource.List)
	r.Route("/memes/{memeId}", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				memeId := chi.URLParam(request, "memeId")
				ctx := context.WithValue(request.Context(), "memeId", memeId)
				next.ServeHTTP(writer, request.WithContext(ctx))
			})
		})
		r.Get("/", memeResource.Get)
		r.Patch("/", memeResource.Update)
		r.Delete("/", memeResource.Delete)
	})
	fs := http.FileServer(config.StaticDir)
	fs = http.StripPrefix(config.StaticPrefix, fs)
	r.Handle(config.StaticPrefix+"*", fs)
	return r
}
