package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := viper.GetString("auth.device_key")
		if r.Header.Get("x-api-key") != key {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type Web struct {
	*gorm.DB
}

func New(db *gorm.DB) *Web {
	return &Web{db}
}

func (w *Web) V1() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/devices", w.devices())
	r.Mount("/assets", w.assets())
	return r
}

func (w *Web) devices() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/ping", func(r chi.Router) {
		r.Use(APIKeyMiddleware)
		r.Post("/", w.PingHandler)
	})
	return r
}

func (w *Web) assets() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(APIKeyMiddleware)
		r.Get("/{id}/download", w.downloadAsset)
	})
	return r
}
