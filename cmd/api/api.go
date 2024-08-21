package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/beng888/buildit/configs"
	"github.com/beng888/buildit/service/auth"
	"github.com/beng888/buildit/service/user"
	"github.com/beng888/buildit/service/web"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {

	sessionStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: configs.Envs.CookiesAuthSecret,
		MaxAge:     configs.Envs.CookiesAuthAgeInSeconds,
		Secure:     configs.Envs.CookiesAuthIsSecure,
		HttpOnly:   configs.Envs.CookiesAuthIsHttpOnly,
	})

	authService := auth.NewAuthService(sessionStore)

	router := http.NewServeMux()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	authHandler := auth.NewHandler(userStore, authService)
	authHandler.RegisterRoutes(router)

	webHandler := web.NewHandler(authService)
	webHandler.RegisterRoutes(router)

	router.Handle("/public/*", disableCache(staticDev()))

	// Serve static files from the public directory
	// fs := http.FileServer(http.Dir("public"))
	// router.Handle("/public/", http.StripPrefix("/public/", fs))

	// subrouter := http.NewServeMux()

	// subrouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// subrouter.HandleFunc("/", templ.Handler(web.HelloForm()))

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

func staticDev() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public")))
}

func disableCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
