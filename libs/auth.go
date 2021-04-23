package libs

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// ValidateSignupForm - Validates user input submitted in the signup form
func ValidateSignupForm(email string, password string, firstName string, lastName string) error {
	if !validateEmail(email) {
		return errors.New("Invalid email")
	}

	if !validatePassword(password) {
		return errors.New("Invalid password")
	}

	if !validateName(firstName) {
		return errors.New("Invalid first name")
	}

	if !validateName(lastName) {
		return errors.New("Invalid last name")
	}

	return nil
}

// HashPassword - Hashed the password and returns it
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ValidateLoginForm - Validates user input submitted in the login form
func ValidateLoginForm(email string, password string, firstName string, lastName string) error {
	if firstName != "" || lastName != "" {
		return errors.New("Invalid request")
	}

	if !validateEmail(email) {
		return errors.New("Invalid email")
	}

	if !validatePassword(password) {
		return errors.New("Invalid password")
	}

	return nil
}

// CheckHash - Compares input password with stored hash
func CheckHash(inputPwd string, hashedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
	if err != nil {
		return false
	}

	return true
}

func validateEmail(email string) bool {
	rxp := regexp.MustCompile("^[a-z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)*$")

	return rxp.MatchString(email)
}

func validatePassword(password string) bool {
	rxp := regexp.MustCompile("^[A-Za-z0-9_@.# &+-]{8,16}$")

	return rxp.MatchString(password)
}

func validateName(name string) bool {
	rxp := regexp.MustCompile("^[A-Za-z]+$")

	return rxp.MatchString(name)
}
