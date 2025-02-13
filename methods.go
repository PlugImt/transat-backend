package main

import (
	"fmt"
	"math/rand"
	"regexp"
)

func checkEmail(email string) (bool, error) {
	regex := `^[a-z0-9.\-]+@imt-atlantique.net$`
	return regexp.MatchString(regex, email)
}

func generate2FACode(digits int) string {
	return fmt.Sprintf("%0"+fmt.Sprint(digits)+"d", rand.Intn(1000000))
}
