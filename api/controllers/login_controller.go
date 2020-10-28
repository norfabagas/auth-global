package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/norfabagas/auth-global/api/jwt"
	"github.com/norfabagas/auth-global/api/models"
	"github.com/norfabagas/auth-global/api/responses"
	"github.com/norfabagas/auth-global/api/utils/crypto"
	"github.com/norfabagas/auth-global/api/utils/formatting"
	"github.com/norfabagas/auth-global/api/utils/smtp"
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

func (server *Server) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()

	keys := r.URL.Query()
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
	err = user.Validate("forget")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userFound, err := user.FindUserByEmail(server.DB, user.Email)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if userFound.ID == 0 {
		responses.ERROR(w, http.StatusNotFound, formatting.FormatError("notFound"))
		return
	}

	generatedPassword := crypto.MD5Hash(time.Now().String())

	changedUser, err := user.ChangePassword(server.DB, userFound.ID, generatedPassword)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	notifyEmail := keys.Get("notify")
	if notifyEmail != "" && notifyEmail == "true" {
		message := fmt.Sprintf("Hello %s,\nWe would like to inform your newly generated password. Please use below:\n%s\n\nThanks", changedUser.Email, generatedPassword)
		err = smtp.Send(
			[]string{changedUser.Email},
			[]string{},
			"Change Password",
			message,
		)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	// if display newly generated password
	visible := keys.Get("visible")
	if visible != "" && visible == "true" {
		responses.JSON(w, http.StatusOK, true, "new password generated", struct {
			Email             string    `json:"email"`
			GeneratedPassword string    `json:"generated_password"`
			RequestTime       time.Time `json:"request_time"`
		}{
			Email:             changedUser.Email,
			GeneratedPassword: generatedPassword,
			RequestTime:       requestTime,
		})
	} else {
		responses.JSON(w, http.StatusOK, true, "Kindly check your email inbox/spam", struct {
			Email       string    `json:"email"`
			RequestTime time.Time `json:"request_time"`
		}{
			Email:       changedUser.Email,
			RequestTime: requestTime,
		})
	}

}

func (server *Server) TestMail(w http.ResponseWriter, r *http.Request) {
	err := smtp.Send([]string{"akunnyagugel@gmail.com"}, []string{}, "Lorem Ipsum", "Lorem Ipsum Dolor Sit Amet")
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, true, "OK", "Success")
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
		return "", formatting.FormatError(err.Error())
	}

	stringID := strconv.Itoa(int(user.ID))

	publicID, err := crypto.Encrypt(stringID, os.Getenv("APP_KEY"))
	if err != nil {
		return "", err
	}
	return jwt.CreateToken(publicID)
}
