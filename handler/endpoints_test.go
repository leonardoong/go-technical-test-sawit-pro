package handler

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestContext(requestBody string) (echo.Context, error) {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, nil
}

func newTestContextWithToken(requestBody string, token string) (echo.Context, error) {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, nil
}

func TestUserRegistration(t *testing.T) {
	repo := new(mocks.RepositoryInterface)

	type args struct {
		requestBody string
	}

	var tests = []struct {
		name   string
		args   args
		mock   func()
		assert func(error, echo.Context)
	}{
		{
			name: "success",
			args: args{
				requestBody: `{"full_name":"John Doe","phone_number":"+628123456789","password":"P@ssw0rd"}`,
			},
			mock: func() {
				repo.On("InsertUser", mock.Anything, mock.Anything).Return(repository.InsertUserOutput{
					UserID: 1,
				}, nil).Once()
			},
			assert: func(err error, ctx echo.Context) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, ctx.Response().Status)
			},
		},
		{
			name: "bad request - body missing",
			args: args{},
			mock: func() {
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
			},
		},
		{
			name: "bad request - partial body missing",
			args: args{
				requestBody: `{"phone_number":"+628123456789","password":"P@ssw0rd"}`,
			},
			mock: func() {
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
			},
		},
		{
			name: "fail - insert db",
			args: args{
				requestBody: `{"full_name":"John Doe","phone_number":"+628123456789","password":"P@ssw0rd"}`,
			},
			mock: func() {
				repo.On("InsertUser", mock.Anything, mock.Anything).Return(repository.InsertUserOutput{}, errors.New("error")).Once()
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
			},
		},
	}

	for _, tt := range tests {
		tt.mock()
		s := Server{
			Repository: repo,
		}

		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := newTestContext(tt.args.requestBody)

			err := s.UserRegistration(ctx)

			tt.assert(err, ctx)
		})
	}
}

func TestLogin(t *testing.T) {
	repo := new(mocks.RepositoryInterface)

	prvKey, err := os.ReadFile("../cert/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}

	pubKey, err := os.ReadFile("../cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	jwtToken := config.NewJWT(prvKey, pubKey)

	type args struct {
		requestBody string
	}

	var tests = []struct {
		name   string
		args   args
		mock   func()
		assert func(error, echo.Context)
	}{
		{
			name: "success",
			args: args{
				requestBody: `{"phone_number":"+6281223129","password":"Leo9999#"}`,
			},
			mock: func() {
				repo.On("GetLoginData", mock.Anything, repository.GetLoginDataInput{
					PhoneNumber: "+6281223129",
				}).Return(repository.GetLoginDataOutput{
					UserID:         1,
					FullName:       "Leonardo",
					HashedPassword: "$2a$04$a2o3BiK8KiH79TFt9QE1hOutA9115oKSUIYQpFAoLldhotz7pwYQe",
				}, nil).Once()

				repo.On("UpdateSuccessfulLogin", mock.Anything, mock.Anything).Return(nil).Once()
			},
			assert: func(err error, ctx echo.Context) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, ctx.Response().Status)
			},
		},
		{
			name: "bad request - body missing",
			args: args{},
			mock: func() {
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
			},
		},
		{
			name: "bad request - partial body missing",
			args: args{
				requestBody: `{"phone_number":"+6281223129"}`,
			},
			mock: func() {
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
			},
		},
		{
			name: "fail - get login data",
			args: args{
				requestBody: `{"phone_number":"+6281223129","password":"Leo9999#"}`,
			},
			mock: func() {
				repo.On("GetLoginData", mock.Anything, repository.GetLoginDataInput{
					PhoneNumber: "+6281223129",
				}).Return(repository.GetLoginDataOutput{}, errors.New("error")).Once()
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
			},
		},
		{
			name: "fail - update successful login",
			args: args{
				requestBody: `{"phone_number":"+6281223129","password":"Leo9999#"}`,
			},
			mock: func() {
				repo.On("GetLoginData", mock.Anything, repository.GetLoginDataInput{
					PhoneNumber: "+6281223129",
				}).Return(repository.GetLoginDataOutput{
					UserID:         1,
					FullName:       "Leonardo",
					HashedPassword: "$2a$04$a2o3BiK8KiH79TFt9QE1hOutA9115oKSUIYQpFAoLldhotz7pwYQe",
				}, nil).Once()

				repo.On("UpdateSuccessfulLogin", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
			},
		},
	}

	for _, tt := range tests {
		tt.mock()
		s := Server{
			Repository: repo,
			Config: &config.Config{
				JWT: jwtToken,
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := newTestContext(tt.args.requestBody)

			err := s.Login(ctx)

			tt.assert(err, ctx)
		})
	}
}

func TestUsers(t *testing.T) {
	repo := new(mocks.RepositoryInterface)

	prvKey, err := os.ReadFile("../cert/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}

	pubKey, err := os.ReadFile("../cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	jwtToken := config.NewJWT(prvKey, pubKey)

	token, _ := jwtToken.Create(time.Minute*1, model.User{
		UserID: 1,
	})

	type args struct {
		token string
	}

	var tests = []struct {
		name   string
		args   args
		mock   func()
		assert func(error, echo.Context)
	}{
		{
			name: "success",
			args: args{
				token: token,
			},
			mock: func() {
				repo.On("GetUserDataByUserID", mock.Anything, repository.GetUserDataByUserIDInput{
					UserID: 1,
				}).Return(model.User{
					UserID:   1,
					FullName: "Leonardo",
				}, nil).Once()

			},
			assert: func(err error, ctx echo.Context) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, ctx.Response().Status)
			},
		},
		{
			name: "token missing",
			args: args{},
			mock: func() {},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusForbidden, ctx.Response().Status)
			},
		},
		{
			name: "fail - get user data",
			args: args{
				token: token,
			},
			mock: func() {
				repo.On("GetUserDataByUserID", mock.Anything, repository.GetUserDataByUserIDInput{
					UserID: 1,
				}).Return(model.User{}, errors.New("error")).Once()

			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
			},
		},
	}

	for _, tt := range tests {
		tt.mock()
		s := Server{
			Repository: repo,
			Config: &config.Config{
				JWT: jwtToken,
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := newTestContextWithToken("", tt.args.token)

			err := s.Users(ctx)

			tt.assert(err, ctx)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	repo := new(mocks.RepositoryInterface)

	prvKey, err := os.ReadFile("../cert/id_rsa")
	if err != nil {
		log.Fatalln(err)
	}

	pubKey, err := os.ReadFile("../cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	jwtToken := config.NewJWT(prvKey, pubKey)

	token, _ := jwtToken.Create(time.Minute*1, model.User{
		UserID: 1,
	})

	type args struct {
		token       string
		requestBody string
	}

	var tests = []struct {
		name   string
		args   args
		mock   func()
		assert func(error, echo.Context)
	}{
		{
			name: "success",
			args: args{
				token:       token,
				requestBody: `{"phone_number":"+6281233245","full_name":"leo"}`,
			},
			mock: func() {
				repo.On("UpdateUserData", mock.Anything, repository.UpdateUserDataInput{
					UserID: 1,
					Data: map[string]string{
						"phone_number": "+6281233245",
						"full_name":    "leo",
					},
				}).Return(nil).Once()
			},
			assert: func(err error, ctx echo.Context) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, ctx.Response().Status)
			},
		},
		{
			name: "token missing",
			args: args{},
			mock: func() {},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusForbidden, ctx.Response().Status)
			},
		},
		{
			name: "fail - update user data",
			args: args{
				token:       token,
				requestBody: `{"phone_number":"+6281233245","full_name":"leo"}`,
			},
			mock: func() {
				repo.On("UpdateUserData", mock.Anything, repository.UpdateUserDataInput{
					UserID: 1,
					Data: map[string]string{
						"phone_number": "+6281233245",
						"full_name":    "leo",
					},
				}).Return(errors.New("error")).Once()
			},
			assert: func(err error, ctx echo.Context) {
				assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
			},
		},
	}

	for _, tt := range tests {
		tt.mock()
		s := Server{
			Repository: repo,
			Config: &config.Config{
				JWT: jwtToken,
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := newTestContextWithToken(tt.args.requestBody, tt.args.token)

			err := s.UpdateUser(ctx)

			tt.assert(err, ctx)
		})
	}
}
