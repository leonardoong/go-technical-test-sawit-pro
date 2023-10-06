package model

import (
	"regexp"
	"strings"
)

type User struct {
	UserID         int32
	FullName       string
	Password       string
	PhoneNumber    string
	SuccesfulLogin int32
}

func (u *User) ValidateRegisterUser() (isValid bool, errorMessages []string) {

	isValidPhoneNumber, errorPhoneNumber := u.ValidatePhoneNumber()
	if !isValidPhoneNumber {
		errorMessages = append(errorMessages, errorPhoneNumber...)
	}

	isValidFullName, errorFullname := u.ValidateFullName()
	if !isValidFullName {
		errorMessages = append(errorMessages, errorFullname...)
	}

	isValidPassword, errorPassword := u.ValidatePassword()
	if !isValidPassword {
		errorMessages = append(errorMessages, errorPassword...)
	}

	return len(errorMessages) == 0, errorMessages
}

func (u *User) ValidatePhoneNumber() (isValid bool, errorMessages []string) {
	// Check phone number length
	if len(u.PhoneNumber) < 10 || len(u.PhoneNumber) > 13 {
		errorMessages = append(errorMessages, "Phone number must be between 10 and 13 characters")
	}

	// Check phone number start with +62
	phoneRegex := `^\+62`
	if !regexp.MustCompile(phoneRegex).MatchString(u.PhoneNumber) {
		errorMessages = append(errorMessages, "Invalid phone number")
	}

	return len(errorMessages) == 0, errorMessages
}

func (u *User) ValidateFullName() (isValid bool, errorMessages []string) {
	// Check full name length
	if len(u.FullName) < 3 || len(u.FullName) > 60 {
		errorMessages = append(errorMessages, "Full name must be between 3 and 60 characters")
	}

	return len(errorMessages) == 0, errorMessages
}

func (u *User) ValidatePassword() (isValid bool, errorMessages []string) {

	// Check the length constraint
	if len(u.Password) < 6 || len(u.Password) > 64 {
		errorMessages = append(errorMessages, "Password must be between 6 and 64 characters")
	}

	// Check for at least one uppercase letter
	if !strings.ContainsAny(u.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		errorMessages = append(errorMessages, "Password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	if !strings.ContainsAny(u.Password, "0123456789") {
		errorMessages = append(errorMessages, "Password must contain at least one digit")
	}

	// Check for at least one special (non-alphanumeric) character
	specialChars := "~!@#$%^&*()-_+=<>?/[]{}|"
	if !strings.ContainsAny(u.Password, specialChars) {
		errorMessages = append(errorMessages, "Password must contain at least one special character")
	}

	return len(errorMessages) == 0, errorMessages
}
