package controllers

import (
	"net/http"

	"github.com/norfabagas/auth-global/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), "OK")
}
