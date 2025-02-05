package main

import (
	"regexp"
)

func checkEmail(email string) (bool, error) {
	regex := `^[a-z0-9.\-]+@imt-atlantique.net$`
	return regexp.MatchString(regex, email)
}
