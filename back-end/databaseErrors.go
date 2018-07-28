package main

import "fmt"

type NotUniqueUserError struct{}

func (e *NotUniqueUserError) Error() string {
	return fmt.Sprintf("Account with that information already exists, try again with a different email or username.")
}

func NewNotUniqueUserError() error {
	return &NotUniqueUserError{}
}
