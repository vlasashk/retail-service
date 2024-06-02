package converter

import (
	"strconv"

	"route256/cart/internal/cart/models"
)

func UserToInt(val string) (int64, error) {
	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil || userID <= 0 {
		return 0, models.ErrInvalidUserID
	}
	return userID, nil
}

func SKUtoInt(val string) (int64, error) {
	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil || userID <= 0 {
		return 0, models.ErrInvalidSKUID
	}
	return userID, nil
}
