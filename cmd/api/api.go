package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/beng888/buildit/service/user"
	"github.com/beng888/buildit/views"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	subrouter := http.NewServeMux()
	subrouter.Handle("/", templ.Handler(views.HelloForm()))

	subrouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// Serve static files from the public directory
	fs := http.FileServer(http.Dir("public"))
	subrouter.Handle("/public/", http.StripPrefix("/public/", fs))

	// subrouter.HandleFunc("/", templ.Handler(web.HelloForm()))

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, subrouter)
}
