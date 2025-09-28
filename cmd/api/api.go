package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Soumil-2007/file-sharing-webApp/services/user"
	"github.com/Soumil-2007/file-sharing-webApp/services/files"
	"github.com/Soumil-2007/file-sharing-webApp/services/middleware"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Users
	userStore := user.NewStore(s.db)
userHandler := user.NewHandler(userStore)
userHandler.RegisterRoutes(apiRouter)

fileStore := files.NewStore(s.db)
fileHandler := files.NewHandler(fileStore)
fileHandler.RegisterRoutes(apiRouter, middleware.AuthMiddleware(userStore))

	// Serve static files (fallback for frontend)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
