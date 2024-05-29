package converter

import (
	"errors"
	"strconv"

	"route256/cart/internal/cart/models/constants"
)

func UserToInt(val string) (int64, error) {
	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil || userID <= 0 {
		return 0, errors.New(constants.ErrInvalidUserID)
	}
	return userID, nil
}

func SKUtoInt(val string) (int64, error) {
	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil || userID <= 0 {
		return 0, errors.New(constants.ErrInvalidSKUID)
	}
	return userID, nil
}
