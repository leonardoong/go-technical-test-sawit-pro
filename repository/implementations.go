package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/SawitProRecruitment/UserService/model"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) InsertUser(ctx context.Context, input InsertUserInput) (output InsertUserOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"INSERT INTO users (phone_number, full_name, password) VALUES ($1, $2, $3) RETURNING id",
		input.PhoneNumber,
		input.FullName,
		input.Password,
	).Scan(&output.UserID)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetLoginData(ctx context.Context, input GetLoginDataInput) (output GetLoginDataOutput, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT id, full_name, password FROM users WHERE phone_number = $1",
		input.PhoneNumber,
	).Scan(&output.UserID, &output.FullName, &output.HashedPassword)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateSuccessfulLogin(ctx context.Context, input UpdateSuccessfulLoginInput) (err error) {
	_, err = r.Db.ExecContext(
		ctx,
		"UPDATE users SET successful_login = successful_login + 1 WHERE phone_number = $1",
		input.PhoneNumber,
	)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUserData(ctx context.Context, input UpdateUserDataInput) (err error) {
	if len(input.Data) == 0 {
		return fmt.Errorf("update user data is empty")
	}

	var (
		query    = "UPDATE users SET %s WHERE id = $1"
		setQuery = []string{}
	)

	if fullName, ok := input.Data["full_name"]; ok {
		setQuery = append(setQuery, fmt.Sprintf("full_name = '%s'", fullName))
	}

	if phoneNumber, ok := input.Data["phone_number"]; ok {
		setQuery = append(setQuery, fmt.Sprintf("phone_number = '%s'", phoneNumber))
	}

	query = fmt.Sprintf(query, strings.Join(setQuery, ","))

	_, err = r.Db.ExecContext(
		ctx,
		query,
		input.UserID,
	)
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetUserDataByUserID(ctx context.Context, input GetUserDataByUserIDInput) (out model.User, err error) {
	err = r.Db.QueryRowContext(
		ctx,
		"SELECT id, full_name, phone_number, successful_login FROM users WHERE id = $1",
		input.UserID,
	).Scan(&out.UserID, &out.FullName, &out.PhoneNumber, &out.SuccesfulLogin)
	if err != nil {
		return
	}
	return
}
