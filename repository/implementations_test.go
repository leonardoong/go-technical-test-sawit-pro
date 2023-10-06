package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var u = &model.User{
	UserID:         1,
	FullName:       "leo",
	Password:       "password",
	PhoneNumber:    "+628123456789",
	SuccesfulLogin: 10,
}

func TestInsertUser(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}

	query := "INSERT INTO users \\(phone_number, full_name, password\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id"

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(u.UserID)

	// test 1 insert success
	mock.ExpectQuery(query).WithArgs(u.PhoneNumber, u.FullName, u.Password).WillReturnRows(rows)

	user, err := repo.InsertUser(context.Background(), InsertUserInput{
		PhoneNumber: u.PhoneNumber,
		FullName:    u.FullName,
		Password:    u.Password,
	})
	assert.NotNil(t, user)
	assert.NoError(t, err)

	// test 2 insert success
	mock.ExpectQuery(query).WithArgs(u.PhoneNumber, u.FullName, u.Password).WillReturnError(sql.ErrConnDone)

	userError, err := repo.InsertUser(context.Background(), InsertUserInput{
		PhoneNumber: u.PhoneNumber,
		FullName:    u.FullName,
		Password:    u.Password,
	})
	assert.Empty(t, userError)
	assert.Error(t, err)
}

func TestGetLoginData(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}

	query := "SELECT id, full_name, password FROM users WHERE phone_number = \\$1"

	rows := sqlmock.NewRows([]string{"id", "full_name", "password"}).
		AddRow(u.UserID, u.FullName, u.Password)

	// test 1 get success
	mock.ExpectQuery(query).WithArgs(u.PhoneNumber).WillReturnRows(rows)

	users, err := repo.GetLoginData(context.Background(), GetLoginDataInput{
		PhoneNumber: u.PhoneNumber,
	})
	assert.NotNil(t, users)
	assert.NoError(t, err)

	// test 2 get error
	mock.ExpectQuery(query).WithArgs(u.PhoneNumber).WillReturnError(sql.ErrConnDone)
	usersError, err := repo.GetLoginData(context.Background(), GetLoginDataInput{
		PhoneNumber: u.PhoneNumber,
	})

	assert.Empty(t, usersError)
	assert.Error(t, err)
}

func TestUpdateSuccessfulLogin(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}

	query := "UPDATE users SET successful_login = successful_login \\+ 1 WHERE phone_number = \\$1"

	// test 1 update success
	mock.ExpectExec(query).WithArgs(u.PhoneNumber).WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateSuccessfulLogin(context.Background(), UpdateSuccessfulLoginInput{
		PhoneNumber: u.PhoneNumber,
	})
	assert.NoError(t, err)

	// test 2 update error
	mock.ExpectExec(query).WithArgs(u.PhoneNumber).WillReturnError(sql.ErrConnDone)

	err = repo.UpdateSuccessfulLogin(context.Background(), UpdateSuccessfulLoginInput{
		PhoneNumber: u.PhoneNumber,
	})
	assert.Error(t, err)
}

func TestGetUserDataByUserID(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}

	query := "SELECT id, full_name, phone_number, successful_login FROM users WHERE id = \\$1"

	rows := sqlmock.NewRows([]string{"id", "full_name", "phone_number", "successful_login"}).
		AddRow(u.UserID, u.FullName, u.PhoneNumber, u.SuccesfulLogin)

	// test 1 get success
	mock.ExpectQuery(query).WithArgs(u.UserID).WillReturnRows(rows)

	users, err := repo.GetUserDataByUserID(context.Background(), GetUserDataByUserIDInput{
		UserID: u.UserID,
	})
	assert.NotNil(t, users)
	assert.NoError(t, err)

	// test 2 get error
	mock.ExpectQuery(query).WithArgs(u.UserID).WillReturnError(sql.ErrConnDone)
	usersError, err := repo.GetUserDataByUserID(context.Background(), GetUserDataByUserIDInput{
		UserID: u.UserID,
	})

	assert.Empty(t, usersError)
	assert.Error(t, err)
}

func TestUpdateUserData(t *testing.T) {
	db, mock := NewMock()
	repo := &Repository{db}

	// test 1 only phone number
	onlyPhoneNumberQuery := fmt.Sprintf("UPDATE users SET phone_number = \\'\\%s\\' WHERE id = \\$1", u.PhoneNumber)
	mock.ExpectExec(onlyPhoneNumberQuery).WithArgs(u.UserID).WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateUserData(context.Background(), UpdateUserDataInput{
		UserID: u.UserID,
		Data: map[string]string{
			"phone_number": u.PhoneNumber,
		},
	})
	assert.NoError(t, err)

	// test 2 only full name
	onlyFullNameQuery := fmt.Sprintf("UPDATE users SET full_name = \\'%s\\' WHERE id = \\$1", u.FullName)
	mock.ExpectExec(onlyFullNameQuery).WithArgs(u.UserID).WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateUserData(context.Background(), UpdateUserDataInput{
		UserID: u.UserID,
		Data: map[string]string{
			"full_name": u.FullName,
		},
	})
	assert.NoError(t, err)

	// test 3 both phone number and full name
	bothQuery := fmt.Sprintf("UPDATE users SET full_name = \\'%s\\'\\,phone_number = \\'\\%s\\' WHERE id = \\$1", u.FullName, u.PhoneNumber)
	mock.ExpectExec(bothQuery).WithArgs(u.UserID).WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateUserData(context.Background(), UpdateUserDataInput{
		UserID: u.UserID,
		Data: map[string]string{
			"full_name":    u.FullName,
			"phone_number": u.PhoneNumber,
		},
	})
	assert.NoError(t, err)

	// test 4 empty data
	err = repo.UpdateUserData(context.Background(), UpdateUserDataInput{
		UserID: u.UserID,
	})
	assert.Error(t, err)

	// test 5 update error
	mock.ExpectExec(bothQuery).WithArgs(u.UserID).WillReturnError(sql.ErrConnDone)
	err = repo.UpdateUserData(context.Background(), UpdateUserDataInput{
		UserID: u.UserID,
		Data: map[string]string{
			"full_name":    u.FullName,
			"phone_number": u.PhoneNumber,
		},
	})
	assert.Error(t, err)
}
