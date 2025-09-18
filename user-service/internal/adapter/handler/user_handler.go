package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/simonaditia/nyayurin/user-service/config"
	"github.com/simonaditia/nyayurin/user-service/internal/adapter"
	"github.com/simonaditia/nyayurin/user-service/internal/adapter/handler/request"
	"github.com/simonaditia/nyayurin/user-service/internal/adapter/handler/response"
	"github.com/simonaditia/nyayurin/user-service/internal/core/domain/entity"
	"github.com/simonaditia/nyayurin/user-service/internal/core/service"
)

type UserHandler interface {
	SignIn(ctx echo.Context) error
}

type userHandler struct {
	userService service.UserServiceInterface
}

var err error

// SignIn implements UserHandler.
func (u *userHandler) SignIn(c echo.Context) error {
	var (
		req        = request.SignInRequest{}
		resp       = response.DefaultResponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] SignIn bind error: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err := c.Validate(&req); err != nil {
		log.Errorf("[UserHandler-2] SignIn validation error: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := u.userService.SigIn(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			log.Warnf("[UserHandler-3] User not found for email %s", req.Email)
			resp.Message = "User not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		log.Errorf("[UserHandler-3] SignIn service error: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	respSignIn.ID = user.ID
	respSignIn.Name = user.Name
	respSignIn.Email = user.Email
	respSignIn.Role = user.RoleName
	respSignIn.Lat = user.Lat
	respSignIn.Lng = user.Lng
	respSignIn.Phone = user.Phone
	respSignIn.AccessToken = token

	resp.Message = "Sign in successful"
	resp.Data = respSignIn

	return c.JSON(http.StatusOK, resp)
}

func NewUserHandler(e *echo.Echo, userService service.UserServiceInterface, cfg *config.Config) UserHandler {
	userHandler := &userHandler{
		userService: userService,
	}

	e.Use(middleware.Recover())
	e.POST("/signin", userHandler.SignIn)

	mid := adapter.NewMiddlewareAdapter(cfg)
	// e.Use(mid.CheckToken())
	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/check", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return userHandler
}
