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

	return errors.New("incorrect details")
}