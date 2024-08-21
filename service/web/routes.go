package web

import (
	"net/http"

	"github.com/beng888/buildit/service/auth"
	"github.com/beng888/buildit/views"
)

type Handler struct {
	auth *auth.AuthService
}

func NewHandler(auth *auth.AuthService) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	// router.HandleFunc("/", auth.RequireAuth(h.handleHome, h.auth))
	router.HandleFunc("/", h.handleHome)
}

func (h *Handler) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.auth.GetSessionUser(r)
	if err != nil {
		// log.Println(err)
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		views.Home().Render(r.Context(), w)
		return
	}

	views.Dashboard(user).Render(r.Context(), w)
}
