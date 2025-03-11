package handler

import (
	"battleship/pkg/common"
	"battleship/pkg/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HandlerAuth struct {
	services *service.Service
}

func NewHandlerAuth(services *service.Service) *HandlerAuth {
	return &HandlerAuth{services: services}
}

func (h *HandlerAuth) InitRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/sign-up", h.signUp)
	router.HandleFunc("/api/sign-in", h.signIn)
}

type ResponseSuccess struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Success  bool   `json:"success"`
	Token    string `json:"access"`
}

func newResponseSuccess(user common.User, message string) ResponseSuccess {
	res := ResponseSuccess{
		Success:  true,
		Message:  message,
		UserID:   user.Id,
		Username: user.Username,
	}
	return res
}

func newResponseSuccessWithToken(user common.User, message, token string) ResponseSuccess {
	res := ResponseSuccess{
		Success:  true,
		Message:  message,
		UserID:   user.Id,
		Username: user.Username,
		Token:    token,
	}
	return res
}

type ResponseError struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func newResponseError(message string) ResponseError {
	res := ResponseError{}
	res.Success = false
	res.Message = message
	return res
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, logMessage string) {
	log.Println(logMessage)
	w.WriteHeader(statusCode)
	response := newResponseError(message)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal("[writeErrorResponse] Failed to encode response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func logStartEnd(handlerName string) func() {
	log.Println("[" + handlerName + "] start")
	return func() {
		log.Println("[" + handlerName + "] end")
	}
}

func (h *HandlerAuth) signUp(w http.ResponseWriter, r *http.Request) {
	defer logStartEnd("signup")()

	var body common.User

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[signup] "+err.Error())
		return
	}

	log.Println("[signup] json decode success")

	if body.Username == "" || body.Password == "" {
		message := "Fields 'username' and 'password' are required"
		writeErrorResponse(w, http.StatusBadRequest, message, "[signup] "+message)
		return
	}

	log.Println("[signup] payload is correct")

	existUser, err := h.services.UserExist(body.Username)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[signup] "+err.Error())
		return
	}

	if existUser.Id != 0 {
		message := "user already exist"
		writeErrorResponse(w, http.StatusConflict, message, "[signup] "+message)
		return
	}

	log.Println("[signup] no conflict")

	userID, err := h.services.Authorization.CreateUser(common.User{
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[signup] "+err.Error())
		return
	}

	log.Println("[signup] user created")

	response := newResponseSuccess(common.User{Id: userID, Username: body.Username}, "User registered successfully")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[signup] "+err.Error())
		return
	}
}

func (h *HandlerAuth) signIn(w http.ResponseWriter, r *http.Request) {
	defer logStartEnd("signin")()

	var body common.User

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[signin] "+err.Error())
		return
	}

	log.Println("[signin] json decode success")

	if body.Username == "" || body.Password == "" {
		message := "Fields 'username' and 'password' are required"
		writeErrorResponse(w, http.StatusBadRequest, message, "[signin] "+message)
		return
	}

	log.Println("[signin] payload is correct")

	user, err := h.services.GetUser(body.Username, body.Password)

	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error(), "[signin] "+err.Error())
		return
	}

	token, err := h.services.GenerateToken(body.Username, body.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error(), "[signin] "+err.Error())
		return
	}

	response := newResponseSuccessWithToken(user, "successfully signed in", token)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error(), "[signup] "+err.Error())
		return
	}
}

func (h *HandlerAuth) getProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting profile...")
}
