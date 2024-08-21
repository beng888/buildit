package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beng888/buildit/types"
	"github.com/beng888/buildit/views"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
	store types.UserStore
	auth  *AuthService
}

func NewHandler(store types.UserStore, auth *AuthService) *Handler {
	return &Handler{store: store, auth: auth}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /auth/{provider}", h.handleProviderLogin)
	router.HandleFunc("GET /auth/{provider}/callback", h.handleAuthCallbackFunction)
	router.HandleFunc("GET /login", h.handleLogin)
	router.HandleFunc("GET /logout", h.handleLogout)

	// admin routes
	// router.HandleFunc("/users/{userID}", auth.WithJWTAuth(h.handleGetUser, h.store)).Methods(http.MethodGet)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	views.Login().Render(r.Context(), w)
}

func (h *Handler) handleProviderLogin(w http.ResponseWriter, r *http.Request) {
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Printf("User already authenticated! %+v", u)
		views.Dashboard(u).Render(r.Context(), w)
	} else {
		provider := r.PathValue("provider")
		if provider == "" {
			http.Error(w, "Provider not specified", http.StatusBadRequest)
			return
		}

		r = gothic.GetContextWithProvider(r, provider)
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) handleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	fmt.Println(user)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	err = h.auth.StoreUserSession(w, r, user)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging out...")

	err := gothic.Logout(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	h.auth.RemoveUserSession(w, r)

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
