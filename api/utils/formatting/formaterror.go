package formatting

import (
	"errors"
	"log"
	"strings"
)

func FormatError(err string) error {
	defer log.Println(err)

	if strings.Contains(err, "name") {
		return errors.New("name already taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("email already taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect email or password")
	}
	if strings.Contains(err, "exists") {
		return errors.New("user already exists")
	}
	if strings.Contains(err, "notFound") {
		return errors.New("user not found")
	}

	return errors.New("incorrect details")
}
