package utils

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHasher(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func StringToId(s string) (uint, error) {
	parseUint, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parseUint), nil
}
