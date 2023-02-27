package infraestructure

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jeffleon/oauth-microservice/pkg/user/aplication"
	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Service aplication.UserService
}

// SignUp
//	@Tags			User
//	@Summary		User Signup
//	@Description	User Signup
//	@Accept			json
//	@Produce		json
//  @Param request body domain.Userdto true "query params"
//	@Success		200		{object}	domain.Userdto
//	@Failure		400		{object}	object{status=string,error=error}
//	@Failure		422		{object}	object{status=string,error=error}
//	@Router			/user/signup [post]
func (h *UserHandler) SignUp(ctx *gin.Context) {
	var requestBody domain.Userdto
	if err := ctx.BindJSON(&requestBody); err != nil {
		logrus.Errorf("CreateUser error making binding error: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, domain.StandardResponse{
			Status: "error",
			Error:  fmt.Sprintf("bad request %s", err),
		})
		return
	}
	user, err := h.Service.SingUp(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.StandardResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.StandardResponse{
		Status:   "success",
		DataType: "object",
		Data:     user,
	})
}

// GetUser
//	@Tags			User
//	@Summary		Get User
//	@Description	Get User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"access token"	extensions(x-example=token)
//	@Success		200				{object}	domain.Userdto
//	@Failure		400				{object}	object{status=string,error=error}
//	@Router			/user [get]
func (h *UserHandler) GetUser(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	user, err := h.Service.GetUser(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.StandardResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.StandardResponse{
		Status:   "success",
		DataType: "object",
		Data:     user,
	})
}

// SignIn
//	@Tags			User
//	@Summary		User Signin
//	@Description	User Signin
//	@Accept			json
//  @Param request body domain.SignIndto true "query params"
//	@Success		200		{object}	domain.Tokendto
//	@Failure		400		{object}	object{status=string,error=error}
//	@Failure		422		{object}	object{status=string,error=error}
//	@Router			/user/signin [post]
func (h *UserHandler) SignIn(ctx *gin.Context) {
	var requestBody domain.SignIndto
	if err := ctx.BindJSON(&requestBody); err != nil {
		logrus.Errorf("SignIn error making binding error: %s", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, domain.StandardResponse{
			Status: "error",
			Error:  fmt.Sprintf("bad request %s", err),
		})
		return
	}
	requestBody.Email = strings.ToLower(requestBody.Email)
	token, err := h.Service.SignIn(&requestBody)
	if err != nil {
		logrus.Errorf("SingIn error service: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, domain.StandardResponse{
			Status: "error",
			Error:  fmt.Sprintf("bad request %s", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.StandardResponse{
		Status:   "success",
		DataType: "object",
		Data:     token,
	})
}

// Refresh Token
//	@Tags			User
//	@Summary		User refresh token
//	@Description	User refresh token
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"access token"	extensions(x-example=token)
//	@Success		200				{object}	domain.Tokendto
//	@Failure		400				{object}	object{status=string,error=error}
//	@Router			/user/refreshToken [post]
func (h *UserHandler) RefreshAccessToken(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	token, err := h.Service.RefreshAccessToken(userID)
	if err != nil {
		logrus.Errorf("Error service refresh access token: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, domain.StandardResponse{
			Status: "error",
			Error:  fmt.Sprintf("bad request %s", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.StandardResponse{
		Status:   "success",
		DataType: "object",
		Data:     token,
	})
}

// Logout
//	@Tags			User
//	@Summary		Logout
//	@Description	Logout
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"access token"	extensions(x-example=token)
//	@Success		200				{object}	string
//	@Failure		400				{object}	object{status=string,error=error}
//	@Router			/user/logout [post]
func (h *UserHandler) Logout(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	token := ctx.Request.Header.Get("Authorization")

	if err := h.Service.Logout(token); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.StandardResponse{
			Status: "error",
			Error:  fmt.Sprintf("bad request %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, domain.StandardResponse{
		Status:   "success",
		DataType: "object",
		Data:     fmt.Sprintf("user %d Logged out", userID),
	})
}
