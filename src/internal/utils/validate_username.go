package utils

import "regexp"

var userNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{2,255}$`)

func ValidateUserName(s string) bool {
	return userNameRegex.MatchString(s)
}
