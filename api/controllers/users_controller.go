package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/norfabagas/auth-global/api/models"
	"github.com/norfabagas/auth-global/api/responses"
	"github.com/norfabagas/auth-global/api/utils/formatting"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("register")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := formatting.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, true, http.StatusText(http.StatusCreated), struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		PublicID  string    `json:"public_id"`
		CreatedAt time.Time `json:"created_at"`
	}{
		Name:      userCreated.Name,
		Email:     userCreated.Email,
		PublicID:  userCreated.PublicID,
		CreatedAt: userCreated.CreatedAt,
	})
}
