package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/norfabagas/auth-global/api/jwt"
	"github.com/norfabagas/auth-global/api/models"
	"github.com/norfabagas/auth-global/api/responses"
	"github.com/norfabagas/auth-global/api/utils/crypto"
	"github.com/norfabagas/auth-global/api/utils/formatting"
	"github.com/norfabagas/auth-global/api/utils/smtp"
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
		PublicID:  user.PublicID,
		CreatedAt: userCreated.CreatedAt,
	})
}

func (server *Server) ShowUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	userID, err := jwt.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.DB.Debug().Model(models.User{}).Where("id = ?", userID).Take(&user).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	name, err := crypto.Decrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), struct {
		PublicID   string    `json:"public_id"`
		Name       string    `json:"name"`
		Email      string    `json:"email"`
		LastUpdate time.Time `json:"last_update"`
	}{
		PublicID:   user.PublicID,
		Name:       name,
		Email:      user.Email,
		LastUpdate: user.UpdatedAt,
	})
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	tokenID, err := jwt.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateUser(server.DB, uint32(tokenID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), struct {
		PublicID    string    `json:"public_id"`
		Email       string    `json:"email"`
		Name        string    `json:"name"`
		LastUpdated time.Time `json:"last_updated"`
	}{
		PublicID:    updatedUser.PublicID,
		Email:       updatedUser.Email,
		Name:        updatedUser.Name,
		LastUpdated: updatedUser.UpdatedAt,
	})
}

func (server *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	json.Unmarshal(body, &user)

	tokenID, err := jwt.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	user.Prepare()
	err = user.Validate("password")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	changedUser, err := user.ChangePassword(server.DB, tokenID, user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	notify := keys.Get("notify")
	if notify != "" && notify == "true" {
		subject := "New Password Change!"
		message := "Your password is successfully changed.\nIf this action is not from you, please contact us."
		go smtp.Send([]string{changedUser.Email}, []string{}, subject, message)
	}

	responses.JSON(w, http.StatusOK, true, "new password created", struct {
		PublicID          string    `json:"public_id"`
		Email             string    `json:"email"`
		PasswordChangedAt time.Time `json:"password_changed_at"`
	}{
		PublicID:          changedUser.PublicID,
		Email:             changedUser.Email,
		PasswordChangedAt: changedUser.UpdatedAt,
	})
}
