package utils

import (
	"errors"
	"fmt"
	"os"
)

const LOGIN_STATUSKEY = "loginStatus"
const LOGIN_USERNAMEKEY = "loginName"

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func MakeDir(name string) (bool, error) {
	if ok, _ := Exists(name); ok {
		return true, nil
	}
	if err := os.Mkdir(name, 0755); err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		return false, err
	}
}
