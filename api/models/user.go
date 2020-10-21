package models

import (
	"errors"
	"html"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/norfabagas/auth-global/api/utils/crypto"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;not null;unique" json:"id"`
	PublicID  string    `gorm:"size:255;not null;unique" json:"public_id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func EscapeAndTrimString(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

func (user *User) BeforeSave() error {
	// hash password
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}

	// generate public_id
	publicID := crypto.MD5Hash(user.Email + user.CreatedAt.String())

	// encrypt name
	encryptedName, err := crypto.Encrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		return err
	}

	// assign created values
	user.Password = string(hashedPassword)
	user.PublicID = publicID
	user.Name = encryptedName

	return nil
}

func (user *User) Prepare() {
	user.ID = 0
	user.Name = EscapeAndTrimString(user.Name)
	user.Email = EscapeAndTrimString(user.Email)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if user.Password == "" {
			return errors.New("required password")
		}
		if user.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("invalid email format")
		}

		return nil

	case "register":
		if user.Name == "" {
			return errors.New("required name")
		}
		if user.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("required email")
		}
		if user.Password == "" {
			return errors.New("required password")
		}
		if passwordLength := len([]rune(user.Password)); passwordLength < 8 {
			return errors.New("password minimum is 8 characters")
		}

		return nil

	case "update":
		if user.Name == "" {
			return errors.New("required name")
		}

		return nil

	case "changePassword":
		if user.Password == "" {
			return errors.New("required password")
		}

		return nil

	default:
		return errors.New("undefined action")
	}
}

func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	err = db.Debug().Create(&user).Error
	if err != nil {
		return &User{}, err
	}

	// decrypt name
	user.Name, _ = crypto.Decrypt(user.Name, os.Getenv("APP_KEY"))

	return user, nil
}

func (user *User) FindUserByID(db *gorm.DB, id uint32) (*User, error) {
	err := db.Debug().Model(&User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &User{}, nil
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("user not found")
	}

	return user, nil
}

func (user *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	err := db.Debug().Model(&User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return &User{}, nil
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("user not found")
	}

	return user, nil
}

func (user *User) UpdateUser(db *gorm.DB, id uint32) (*User, error) {
	encryptedName, err := crypto.Encrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		return &User{}, err
	}

	db = db.Debug().Model(&User{}).Where("id = ?").UpdateColumns(
		map[string]interface{}{
			"name":       encryptedName,
			"updated_at": user.UpdatedAt,
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?").Take(&user).Error
	if err != nil {
		return &User{}, err
	}

	user.Name, err = crypto.Decrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) ChangePassword(db *gorm.DB, id uint32, password string) (*User, error) {
	hashedPassword, err := Hash(password)
	if err != nil {
		return &User{}, err
	}

	db = db.Debug().Model(&User{}).Where("id = ?", id).UpdateColumns(
		map[string]interface{}{
			"password": string(hashedPassword),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) DeleteUser(db *gorm.DB, id uint32) (string, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", id).Take(&user).Delete(&user)
	if db.Error != nil {
		return "", db.Error
	}

	return user.Email, nil
}
