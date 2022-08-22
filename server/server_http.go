package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/zufzuf/cake-store/db"
	"github.com/zufzuf/cake-store/handler"
	"github.com/zufzuf/cake-store/libs/logger"
	"github.com/zufzuf/cake-store/libs/util"
	"github.com/zufzuf/cake-store/repository"
	AppMiddleware "github.com/zufzuf/cake-store/server/middleware"
	"github.com/zufzuf/cake-store/service"
)

type CakeHandler interface {
	FindCake(rw http.ResponseWriter, r *http.Request)
	FindAllCake(rw http.ResponseWriter, r *http.Request)
	AddCake(rw http.ResponseWriter, r *http.Request)
	UpdateCake(rw http.ResponseWriter, r *http.Request)
	DeleteCake(rw http.ResponseWriter, r *http.Request)
}

func NewHTTPServer() *HTTPServer {
	db := db.Init()

	logger.StartLogger()
	util.NewValidator()

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{"*"},
		AllowCredentials:   true,
		MaxAge:             60,
		OptionsPassthrough: false,
		Debug:              false,
	}))
	r.Use(AppMiddleware.Tracker)

	repoCake := &repository.Cake{
		DB: db,
	}

	srv := &service.Cake{
		Repo: repoCake,
	}

	server := &HTTPServer{
		Router:      r,
		DB:          db,
		CakeHandler: &handler.Cake{Service: srv},
	}

	server.routes()

	return server
}

type HTTPServer struct {
	Router *chi.Mux
	DB     *sql.DB

	CakeHandler CakeHandler
}

func (hs *HTTPServer) Run(ctx context.Context) error {
	port, ok := os.LookupEnv("API_PORT")
	if !ok {
		port = "3000"
	}

	server := http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           hs.Router,
		IdleTimeout:       0,
		WriteTimeout:      5 * time.Second,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		log.Printf("start cake api")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start / shutdown cake api, err : \n%+v\n", err)
		}
	}()

	<-ctx.Done()

	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("server shutdown")

	if err := server.Shutdown(shutdown); err != nil {
		log.Fatalf("shutdown cake api, err : \n%+v\n", err)
	}

	log.Printf("server shutdown properly")

	if err := hs.DB.Close(); err != nil {
		log.Fatal("unable close db connection")
	}

	return nil
}
