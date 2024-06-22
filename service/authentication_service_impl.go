package service

import (
	"errors"
	"gins/config"
	"gins/data/request"
	"gins/helper"
	"gins/model"
	"gins/repository"
	"gins/utils"

	"github.com/go-playground/validator/v10"
)

type AuthenticationServiceImpl struct {
	UsersRepository repository.UsersRepository
	Validate        *validator.Validate
}

func NewAuthenticationServiceImpl(usersRepository repository.UsersRepository, validate *validator.Validate) AuthenticationService {
	return &AuthenticationServiceImpl{
		UsersRepository: usersRepository,
		Validate:        validate,
	}
}

// Login implements AuthenticationService.
func (a *AuthenticationServiceImpl) Login(users request.LoginRequest) (string, error) {
	// Find username in database

	new_user, user_err := a.UsersRepository.FindByUsername(users.Username)
	if user_err != nil {
		return "", errors.New("invalid username or password")
	}

	config, _ := config.LoadConfig(".")

	verify_error := utils.VerifyPasswrod(new_user.Password, users.Password)
	if verify_error != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate token
	token, err_token := utils.GenerateToken(config.TokenExpiresIn, new_user.Id, config.TokenSecret)
	helper.ErrorPanic(err_token)

	return token, nil
}

// Register implements AuthenticationService.
func (a *AuthenticationServiceImpl) Register(users request.CreateUsersRequest) {
	hashedPassword, err := utils.HashPassword(users.Password)
	helper.ErrorPanic(err)

	newUser := model.Users{
		Username: users.Username,
		Email:    users.Email,
		Password: hashedPassword,
	}

	a.UsersRepository.Save(newUser)
}
