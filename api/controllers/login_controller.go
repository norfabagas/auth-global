package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/norfabagas/auth-global/api/jwt"
	"github.com/norfabagas/auth-global/api/models"
	"github.com/norfabagas/auth-global/api/responses"
	"github.com/norfabagas/auth-global/api/utils/crypto"
	"github.com/norfabagas/auth-global/api/utils/formatting"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
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
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.signIn(user.Email, user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.DB.Debug().Model(models.User{}).Take(&user).Error
	if err != nil {
		formattedError := formatting.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	// decrypt name
	user.Name, err = crypto.Decrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, true, http.StatusText(http.StatusOK), struct {
		Token string `json:"token"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		Token: token,
		Email: user.Email,
		Name:  user.Name,
	})
}

func (server *Server) signIn(email, password string) (string, error) {
	var err error

	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return jwt.CreateToken(user.PublicID)
}
