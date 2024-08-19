package handler

import (
	_ "auth_service/api/docs"
	"auth_service/api/email"
	"auth_service/api/models"
	"auth_service/api/token"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register godoc
// @Router       /auth_service/register [POST]
// @Summary      Register a new user
// @Description  Register a new user with an optional profile image
// @Tags         Auth_Service
// @Accept       json
// @Produce      json
// @Param        user body models.Register true "User information"
// @Success      201  {object}  models.RegisterResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) Register(ctx *gin.Context) {
	var (
		request  models.Register
		response *models.RegisterResponse
		err      error
	)

	//file, header, err := ctx.Request.FormFile("file")
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
	//	return
	//}
	//defer func(file multipart.File) {
	//	if err := file.Close(); err != nil {
	//		h.log.Error("Failed to close file: %v", logger.Error(err))
	//	}
	//}(file)
	//
	//firebaseStorage, err := firebase.NewFirebaseStorage("./configs/service-account.json", "testproject-fc4d3.appspot.com")
	//if err != nil {
	//	h.log.Error("Failed to initialize Firebase Storage: %v", logger.Error(err))
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firebase Storage"})
	//	return
	//}
	//
	//fileURL, err := firebaseStorage.UploadFile(context.Background(), file, header.Filename)
	//if err != nil {
	//	h.log.Error("Failed to upload file: %v", logger.Error(err))
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
	//	return
	//}
	//h.log.Info("File URL: %s")

	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "Error reading request body", http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("++++++", request)
	response, err = h.service.RegisterUser(context.Background(), request)
	if err != nil {
		handleResponse(ctx, h.log, "Error creating user", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(ctx, h.log, "Success", http.StatusCreated, response)
}

// Login godoc
// @Router       /auth_service/login/ [POST]
// @Summary      LoginUser User
// @Description  LoginUser User
// @Tags         Auth_Service
// @Accept       json
// @Produce      json
// @Param        user body models.LoginRequest false "user"
// @Success      200  {object}  models.Token
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) Login(ctx *gin.Context) {
	var (
		err      error
		request  models.LoginRequest
		response *models.LoginResponse
	)
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while reading by body  login of user", http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("sssss", request)
	response, err = h.service.Login(context.Background(), &request)
	if err != nil {
		handleResponse(ctx, h.log, "error is while login  user", http.StatusInternalServerError, err.Error())
		return
	}

	generateJWTToken, err := token.GenerateJWTToken(h.log, response)
	if err != nil {
		handleResponse(ctx, h.log, "error is while generate token for   user", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(ctx, h.log, "", http.StatusOK, generateJWTToken)

}

// Refresh godoc
// @Router       /auth_service/refresh/ [GET]
// @Summary      new_generate_token for User
// @Description  new_generate_token  for User
// @Tags         Auth_Service
// @Accept       json
// @Produce      json
// @Param        token query string false "token "
// @Success      200  {object}  models.Token
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) Refresh(ctx *gin.Context) {
	var (
		err      error
		tokenStr string
		response *models.LoginResponse
	)

	tokenStr = ctx.Query("token")

	response, err = token.NewGenerateToken(h.log, tokenStr)
	if err != nil {
		handleResponse(ctx, h.log, "Error generating new token", http.StatusBadRequest, err.Error())
		return
	}

	newToken, err := token.GenerateJWTToken(h.log, response)
	if err != nil {
		handleResponse(ctx, h.log, "Error generating JWT token", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(ctx, h.log, "Token refreshed successfully", http.StatusOK, newToken)
}

// Logout godoc
// @Router       /auth_service/logout/ [GET]
// @Summary      Logout of  User
// @Description  Logout of User
// @Tags         Auth_Service
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) Logout(ctx *gin.Context) {
	handleResponse(ctx, h.log, "logout successfully", http.StatusOK, "logout successfully")
	return
}

// SendPasswordToEmail godoc
// @Router       /auth_service/reset_password/ [PUT]
// @Summary      Change_Password of  User
// @Description  Change_Password of User
// @Tags         Auth_Service
// @Accept       json
// @Produce      json
// @Param        user query models.SendEmail false "user"
// @Success      200  {object}  string
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) SendPasswordToEmail(ctx *gin.Context) {
	var (
		//err     error
		request models.SendEmail
	)

	if h == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	request.Email = ctx.Query("email")
	request.NewPassword = ctx.Query("new_password")
	request.ConfirmNewPassword = ctx.Query("confirm_new_password")
	email.Email(request.Email)

}
