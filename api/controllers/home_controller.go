package controllers

import (
	"errors"
	"net/http"
	"os"

	"github.com/norfabagas/auth-global/api/jwt"
	"github.com/norfabagas/auth-global/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), "OK")
}

func (server *Server) ApiSecret(w http.ResponseWriter, r *http.Request) {
	token := jwt.ExtractToken(r)
	if token != os.Getenv("ACCEPTED_TOKEN") {
		responses.ERROR(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), os.Getenv("API_SECRET"))
}
