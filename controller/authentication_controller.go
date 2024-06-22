package controller

import (
	"gins/data/request"
	"gins/data/response"
	"gins/helper"
	"gins/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	authenticationService service.AuthenticationService
}

func NewAuthenticationController(service service.AuthenticationService) *AuthenticationController {
	return &AuthenticationController{
		authenticationService: service,
	}
}

func (controller *AuthenticationController) Login(ctx *gin.Context) {
	LoginRequest := request.LoginRequest{}
	err := ctx.ShouldBindJSON(&LoginRequest)
	helper.ErrorPanic(err)

	token, err := controller.authenticationService.Login(LoginRequest)

	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid username or password",
		}

		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	resp := response.LoginResponse{
		TokenType: "Bearer",
		Token:     token,
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "OK",
		Message: "Successfully log in!",
		Data:    resp,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *AuthenticationController) Register(ctx *gin.Context) {
	createUsersRequest := request.CreateUsersRequest{}
	err := ctx.ShouldBindJSON(&createUsersRequest)
	helper.ErrorPanic(err)

	controller.authenticationService.Register(createUsersRequest)
	webResponse := response.Response{
		Code:    200,
		Status:  "OK",
		Message: "Successfully created user!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
