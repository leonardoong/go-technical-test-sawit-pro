// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/model"
)

type RepositoryInterface interface {
	GetLoginData(ctx context.Context, input GetLoginDataInput) (output GetLoginDataOutput, err error)
	GetUserDataByUserID(ctx context.Context, input GetUserDataByUserIDInput) (model.User, error)

	InsertUser(ctx context.Context, in InsertUserInput) (out InsertUserOutput, err error)

	UpdateSuccessfulLogin(ctx context.Context, in UpdateSuccessfulLoginInput) error
	UpdateUserData(ctx context.Context, in UpdateUserDataInput) error
}
