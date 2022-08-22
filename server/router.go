package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (hs *HTTPServer) routes() {
	hs.Router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("🍰 cake api 🔥"))
	})
	hs.Router.Route("/cakes", func(r chi.Router) {
		r.Get("/", hs.CakeHandler.FindAllCake)
		r.Post("/", hs.CakeHandler.AddCake)
		r.Get("/{id:[0-9]+}", hs.CakeHandler.FindCake)
		r.Patch("/{id:[0-9]+}", hs.CakeHandler.UpdateCake)
		r.Delete("/{id:[0-9]+}", hs.CakeHandler.DeleteCake)
	})
}
