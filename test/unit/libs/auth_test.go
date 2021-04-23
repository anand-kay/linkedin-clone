package libs

import (
	"testing"

	"github.com/anand-kay/linkedin-clone/libs"
)

func TestValidateSignupForm(t *testing.T) {
	err := libs.ValidateSignupForm("first@first.com", "sfh43kjfgsd", "John", "Smith")
	if err != nil {
		t.Error("ValidateSignupForm - Failed")
	}

	err = libs.ValidateSignupForm("firstfirst.com", "sfh43kjfgsd", "John", "Smith")
	if err == nil {
		t.Error("ValidateSignupForm - Failed")
	}

	err = libs.ValidateSignupForm("first@first.com", "43gf", "John", "Smith")
	if err == nil {
		t.Error("ValidateSignupForm - Failed")
	}

	err = libs.ValidateSignupForm("first@first.com", "sfh43kjfgsd", "364", "Smith")
	if err == nil {
		t.Error("ValidateSignupForm - Failed")
	}

	err = libs.ValidateSignupForm("first@first.com", "sfh43kjfgsd", "John", "&&&")
	if err == nil {
		t.Error("ValidateSignupForm - Failed")
	}
}

func TestHashPassword(t *testing.T) {
	hashedPwd, err := libs.HashPassword("qwertyuiop")
	if err != nil {
		t.Error("HashPassword - Failed")
	}
	if hashedPwd == "" {
		t.Error("HashPassword - Failed")
	}
}

func TestValidateLoginForm(t *testing.T) {
	err := libs.ValidateLoginForm("first@first.com", "jdkfgh8432nj", "", "")
	if err != nil {
		t.Error("ValidateLoginForm - Failed")
	}

	err = libs.ValidateLoginForm("first@first.com", "jdkfgh8432nj", "Jake", "")
	if err == nil {
		t.Error("ValidateLoginForm - Failed")
	}

	err = libs.ValidateLoginForm("first@first.com", "jdkfgh8432nj", "", "Sparrow")
	if err == nil {
		t.Error("ValidateLoginForm - Failed")
	}

	err = libs.ValidateLoginForm("firstfirst.com", "jdkfgh8432nj", "", "")
	if err == nil {
		t.Error("ValidateLoginForm - Failed")
	}

	err = libs.ValidateLoginForm("first@first.com", "fsgb", "", "")
	if err == nil {
		t.Error("ValidateLoginForm - Failed")
	}
}

func TestCheckHash(t *testing.T) {
	password := "zxcvb12nm"

	hashedPwd, err := libs.HashPassword(password)
	if err != nil {
		t.Error("Error while hashing password")
	}

	isValid := libs.CheckHash(password, hashedPwd)
	if !isValid {
		t.Error("CheckHash - Failed")
	}

	isValid = libs.CheckHash("kjghs94sd", hashedPwd)
	if isValid {
		t.Error("CheckHash - Failed")
	}
}
