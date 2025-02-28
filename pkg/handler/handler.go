package handler

import (
	"battleship/pkg/service"
	"fmt"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/sign-up", h.signUp)
	router.HandleFunc("/api/sign-in", h.signIn)

	return router
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registrration...")
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authentification...")
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting profile...")
}
