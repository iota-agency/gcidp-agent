package utils

import (
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const (
	lowerCharSet    = "abcdedfghijklmnopqrst"
	upperCharSet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet  = "!@#$%&*"
	numberSet       = "0123456789"
	alphaNumericSet = lowerCharSet + upperCharSet + numberSet
	allCharSet      = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func RunCmd(cmd *exec.Cmd) error {
	var outbuf, errbuf strings.Builder
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		return errors.New(errbuf.String())
	}
	return nil
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func MkDirIfNone(dir string) error {
	if DirExists(dir) {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func RandStringFullCharSet(length int) string {
	var password string
	for i := 0; i < length; i++ {
		random := rand.Intn(len(allCharSet))
		password += string(allCharSet[random])
	}
	return password
}

func RandStringLowerCharSet(length int) string {
	var password string
	for i := 0; i < length; i++ {
		random := rand.Intn(len(lowerCharSet))
		password += string(lowerCharSet[random])
	}
	return password
}

func RandStringAlphaNumeric(length int) string {
	var password string
	for i := 0; i < length; i++ {
		random := rand.Intn(len(alphaNumericSet))
		password += string(alphaNumericSet[random])
	}
	return password
}

func RandStringNumeric(length int) string {
	var password string
	for i := 0; i < length; i++ {
		random := rand.Intn(len(numberSet))
		password += string(numberSet[random])
	}
	return password
}

func RandStringSpecial(length int) string {
	var password string
	for i := 0; i < length; i++ {
		random := rand.Intn(len(specialCharSet))
		password += string(specialCharSet[random])
	}
	return password
}
