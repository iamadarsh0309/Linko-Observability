package main

import (
	pkgerr "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// server keeps password validation in the root package so its stack frame is
// visible to the lesson's stack-trace check.
type server struct{}

func (s *server) validatePassword(password, stored string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, pkgerr.WithStack(err)
	}
	return true, nil
}
