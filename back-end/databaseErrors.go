package main

import "fmt"

type NotUniqueGameError struct{}

func (e *NotUniqueGameError) Error() string {
	return fmt.Sprintf("Game with that name already exists, try again with a different email or username.")
}

func NewNotUniqueGameError() error {
	return &NotUniqueGameError{}
}

type NotUniqueUserError struct{}

func (e *NotUniqueUserError) Error() string {
	return fmt.Sprintf("Account with that information already exists, try again with a different email or username.")
}

func NewNotUniqueUserError() error {
	return &NotUniqueUserError{}
}
