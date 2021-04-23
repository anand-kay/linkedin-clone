package libs

import (
	"errors"
	"strconv"

	"github.com/anand-kay/linkedin-clone/utils"
)

// ValidateQueryParams - Validates query params of FetchAllPosts
func ValidateQueryParams(otherUserID string, page string, limit string) error {
	err := utils.ValidateID(otherUserID)
	if err != nil {
		return errors.New("Invalid user id")
	}

	err = validatePage(page)
	if err != nil {
		return err
	}

	err = validateLimit(limit)
	if err != nil {
		return err
	}

	return nil
}

func validatePage(page string) error {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 0 {
		return errors.New("Invalid page")
	}

	return nil
}

func validateLimit(limit string) error {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return errors.New("Invalid limit")
	}

	return nil
}
