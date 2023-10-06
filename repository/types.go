// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type InsertUserInput struct {
	PhoneNumber string
	FullName    string
	Password    string
}

type InsertUserOutput struct {
	UserID int32
}

type GetLoginDataInput struct {
	PhoneNumber string
}

type GetLoginDataOutput struct {
	UserID         int32
	FullName       string
	HashedPassword string
}

type UpdateSuccessfulLoginInput struct {
	PhoneNumber string
}

type UpdateUserDataInput struct {
	UserID int32
	Data   map[string]string
}

type GetUserDataByUserIDInput struct {
	UserID int32
}
