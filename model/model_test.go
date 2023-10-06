package model

import (
	"testing"
)

func TestValidateRegisterUser(t *testing.T) {
	tests := []struct {
		input    User
		expected bool
	}{
		{User{
			FullName:    "Leonardo",
			Password:    "P@ssw0rd",
			PhoneNumber: "+628123456789",
		}, true},
		{User{
			FullName:    "Leonardo",
			Password:    "P@ssw0rd",
			PhoneNumber: "+8123456789",
		}, false}, // Invalid phone number
		{User{
			FullName:    "L",
			Password:    "P@ssw0rd",
			PhoneNumber: "+628123456789",
		}, false}, // Invalid full name
		{User{
			FullName:    "Leonardo",
			Password:    "w0rd",
			PhoneNumber: "+628123456789",
		}, false}, // Invalid password
	}

	for _, test := range tests {
		user := test.input
		isValid, _ := user.ValidateRegisterUser()
		if isValid != test.expected {
			t.Errorf("For input '%+v', expected validation result %v, but got %v", test.input, test.expected, isValid)
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+628123456789", true},
		{"+62123", false},         // Missing digit
		{"+6281234567890", false}, // Too long
		{"628123456789", false},   // Missing '+'
		{"08123456789", false},    // Missing '+62'
		{"+62abc", false},         // Invalid characters
	}

	for _, test := range tests {
		user := User{PhoneNumber: test.input}
		isValid, _ := user.ValidatePhoneNumber()
		if isValid != test.expected {
			t.Errorf("For input '%s', expected validation result %v, but got %v", test.input, test.expected, isValid)
		}
	}
}

func TestValidateFullName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Leonardo", true},
		{"A", false}, // Too short
		{"Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum Lorem Ipsum", false}, // Too long
	}

	for _, test := range tests {
		user := User{FullName: test.input}
		isValid, _ := user.ValidateFullName()
		if isValid != test.expected {
			t.Errorf("For input '%s', expected validation result %v, but got %v", test.input, test.expected, isValid)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"P@ssw0rd", true},
		{"abc123", false}, // Missing uppercase letter
		{"ABCXYZ", false}, // Missing digit
		{"P@ss", false},   // Too short
		{"P@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rdasd", false}, // Too long
	}

	for _, test := range tests {
		user := User{Password: test.input}
		isValid, _ := user.ValidatePassword()
		if isValid != test.expected {
			t.Errorf("For input '%s', expected validation result %v, but got %v", test.input, test.expected, isValid)
		}
	}
}
