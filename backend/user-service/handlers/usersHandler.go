package handlers

import (
	"12factorapp/data"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type KeyUser struct{}

type UsersHandler struct {
	logger   *log.Logger
	userRepo data.UsersRepo
}

// Injecting the logger makes this code much more testable.
func New(l *log.Logger, pr data.UsersRepo) *UsersHandler {
	return &UsersHandler{l, pr}
}

func (u *UsersHandler) GetUsers(rw http.ResponseWriter, h *http.Request) {
	users := u.userRepo.GetUsers()
	err := users.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		u.logger.Println("Unable to convert to json :", err)
		return
	}
}

func (u *UsersHandler) GetAllUsers(rw http.ResponseWriter, h *http.Request) {
	allUsers := u.userRepo.GetAll()
	err := allUsers.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		u.logger.Println("Unable to convert to json :", err)
		return
	}
}

func (u *UsersHandler) PostUsers(rw http.ResponseWriter, h *http.Request) {
	prod := h.Context().Value(KeyUser{}).(*data.User)
	u.userRepo.AddUser(prod)
	rw.WriteHeader(http.StatusCreated)
}

func (u *UsersHandler) PutUsers(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id, _ := strconv.Atoi(vars["id"])

	prod := h.Context().Value(KeyUser{}).(*data.User)
	putErr := u.userRepo.PutUser(prod, id)

	if putErr != nil {
		http.Error(rw, putErr.Error(), http.StatusBadRequest)
		u.logger.Println(putErr.Error())
		return
	}

	err := prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		u.logger.Println("Unable to convert to json :", err)
		return
	}
}

func (u *UsersHandler) DeleteUsers(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		u.logger.Println("Unable to convert from ascii to integer - input was :", vars["id"])
		return
	}

	err = u.userRepo.DeleteUser(id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		u.logger.Println("Unable to delete user.", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

//Middleware to try and decode the incoming body. When decoded we run the validation on it just to check if everything is okay
//with the model. If anything is wrong we terminate the execution and the code won't even hit the handler functions.
//With a key we bind what we read to the context of the current request. Later we use that key to get to the read value.

func (u *UsersHandler) MiddlewareUserValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		prod := &data.User{}
		err := prod.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			u.logger.Println(err)
			return
		}

		err = prod.Validate()

		if err != nil {
			u.logger.Println("Error validating user", err)
			http.Error(rw, fmt.Sprintf("Error validating user: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(h.Context(), KeyUser{}, prod)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

//Middleware to centralize general logging and to add the header values for all requests.

func (u *UsersHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		u.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
