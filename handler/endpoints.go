package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (s *Server) UserRegistration(ctx echo.Context) error {

	var (
		resp          generated.RegisterUserResponse
		errResp       = generated.ErrorResponse{}
		errorMessages []string
	)

	// Get request body data
	body := new(generated.RegisterUserRequest)
	if err := ctx.Bind(body); err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	// Validate request body content exist
	if body.FullName == "" || body.PhoneNumber == "" || body.Password == "" {
		errResp.Message = "Phone Number or Full Name or Password is missing."
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	user := model.User{
		FullName:    body.FullName,
		Password:    body.Password,
		PhoneNumber: body.PhoneNumber,
	}

	isValid, errorMessages := user.ValidateRegisterUser()

	if !isValid || len(errorMessages) > 0 {
		errResp := generated.ErrorResponse{
			Message:       "Invalid Request. Please meet the criteria",
			ErrorMessages: &errorMessages,
		}
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	// Hashed password
	hashedPassword, err := HashedPassword(body.Password)
	if err != nil {
		errResp := generated.ErrorResponse{
			Message: err.Error(),
		}

		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	// Insert user data to DB
	out, err := s.Repository.InsertUser(ctx.Request().Context(), repository.InsertUserInput{
		PhoneNumber: body.PhoneNumber,
		FullName:    body.FullName,
		Password:    hashedPassword,
	})
	if err != nil {
		errResp := generated.ErrorResponse{
			Message: err.Error(),
		}

		// Handle phone number already exist using psql unique constraint
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				errResp.Message = "Phone number already registered. Please use another phone number."
				return ctx.JSON(http.StatusConflict, errResp)
			}
		}

		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	resp.Message = fmt.Sprintf("Successfuly create user with id : %d", out.UserID)

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {

	var (
		resp    generated.LoginResponse
		errResp = generated.ErrorResponse{}
	)

	// Get request body data
	body := new(generated.LoginRequest)
	if err := ctx.Bind(body); err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	// Validate request body content exist
	if body.PhoneNumber == "" || body.Password == "" {
		errResp.Message = "Phone Number or Password is missing."
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	userData, err := s.Repository.GetLoginData(ctx.Request().Context(), repository.GetLoginDataInput{
		PhoneNumber: body.PhoneNumber,
	})
	if err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	// Validate password match
	if !CompareHashAndPassword(userData.HashedPassword, body.Password) {
		errResp.Message = "Invalid phone number or password."
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	// Create JWT token
	token, err := s.Config.JWT.Create(time.Hour*1, model.User{
		UserID: userData.UserID,
	})
	if err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	// Increment succesfull login in database
	err = s.Repository.UpdateSuccessfulLogin(ctx.Request().Context(), repository.UpdateSuccessfulLoginInput{
		PhoneNumber: body.PhoneNumber,
	})
	if err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	resp.Message = fmt.Sprintf("Successfuly login user with id : %d", userData.UserID)
	resp.Jwt = token

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Users(ctx echo.Context) error {

	var (
		resp    generated.UsersResponse
		errResp = generated.ErrorResponse{}
	)

	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		errResp.Message = "Forbidden Code"
		return ctx.JSON(http.StatusForbidden, errResp)
	}

	userData, err := s.Config.JWT.Validate(token)
	if err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusForbidden, errResp)
	}

	out, err := s.Repository.GetUserDataByUserID(ctx.Request().Context(), repository.GetUserDataByUserIDInput{
		UserID: userData.UserID,
	})
	if err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	resp.FullName = out.FullName
	resp.PhoneNumber = out.PhoneNumber

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateUser(ctx echo.Context) error {

	var (
		resp    generated.UpdateUserResponse
		errResp = generated.ErrorResponse{}
	)

	// Get token from request
	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		errResp.Message = "Forbidden Code"
		return ctx.JSON(http.StatusForbidden, errResp)
	}

	// Get request body data
	body := new(generated.UpdateUserRequest)
	if err := ctx.Bind(body); err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	// Validate request body content exist
	if body.FullName == nil && body.PhoneNumber == nil {
		errResp.Message = "Phone Number and Full Name is missing. Nothing to update"
		return ctx.JSON(http.StatusBadRequest, errResp)
	}

	user := model.User{}

	if body.PhoneNumber != nil {
		user.PhoneNumber = *body.PhoneNumber
		validPhoneNumber, errorMessages := user.ValidatePhoneNumber()
		if !validPhoneNumber {
			errResp.Message = strings.Join(errorMessages, " & ")
			return ctx.JSON(http.StatusBadRequest, errResp)
		}
	}

	if body.FullName != nil {
		user.FullName = *body.FullName
		validFullName, errorMessages := user.ValidateFullName()
		if !validFullName {
			errResp.Message = strings.Join(errorMessages, " & ")
			return ctx.JSON(http.StatusBadRequest, errResp)
		}
	}

	// Validate token
	userData, err := s.Config.JWT.Validate(token)
	if err != nil {
		errResp.Message = err.Error()
		return ctx.JSON(http.StatusForbidden, errResp)
	}

	dataToUpdate := make(map[string]string)
	if body.FullName != nil {
		dataToUpdate["full_name"] = user.FullName
	}

	if body.PhoneNumber != nil {
		dataToUpdate["phone_number"] = user.PhoneNumber
	}

	err = s.Repository.UpdateUserData(ctx.Request().Context(), repository.UpdateUserDataInput{
		UserID: userData.UserID,
		Data:   dataToUpdate,
	})
	if err != nil {
		errResp := generated.ErrorResponse{
			Message: err.Error(),
		}

		// Handle phone number already exist using psql unique constraint
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				errResp.Message = "Phone number already registered. Please use another phone number."
				return ctx.JSON(http.StatusConflict, errResp)
			}
		}

		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	resp.Message = "Successfuly update user data."
	return ctx.JSON(http.StatusOK, resp)

}
