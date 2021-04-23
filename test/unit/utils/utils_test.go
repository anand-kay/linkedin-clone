package utils

import (
	"testing"

	"github.com/anand-kay/linkedin-clone/utils"
)

func TestGenerateJWT(t *testing.T) {
	token, err := utils.GenerateJWT(1, "first@first.com", "John", "Smith")
	if err != nil {
		t.Error("GenerateJWT - Failed")
	}
	if token == "" {
		t.Error("GenerateJWT - Failed")
	}
}

func TestValidateID(t *testing.T) {
	err := utils.ValidateID("1")
	if err != nil {
		t.Error("ValidateID - Failed")
	}

	err = utils.ValidateID("1999")
	if err != nil {
		t.Error("ValidateID - Failed")
	}

	err = utils.ValidateID("0")
	if err == nil {
		t.Error("ValidateID - Failed")
	}

	err = utils.ValidateID("ns")
	if err == nil {
		t.Error("ValidateID - Failed")
	}

	err = utils.ValidateID("-1")
	if err == nil {
		t.Error("ValidateID - Failed")
	}
}
