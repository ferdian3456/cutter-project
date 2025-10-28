package http

import (
	"cutterproject/internal/constant"
	"cutterproject/internal/model"
	"cutterproject/internal/usecase"
	"cutterproject/internal/util"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type UserController struct {
	UserUsecase *usecase.UserUsecase
	Log         *zap.Logger
	Config      *koanf.Koanf
}

func NewUserController(userUsecase *usecase.UserUsecase, zap *zap.Logger, koanf *koanf.Koanf) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
		Log:         zap,
		Config:      koanf,
	}
}

func (controller UserController) Register(ctx *fiber.Ctx) error {
	var payload model.UserCreateRequest
	err := util.ReadRequestBody(ctx, &payload)
	if err != nil {
		return util.SendErrorResponse(ctx, &model.ValidationError{
			Code:    constant.ERR_INVALID_REQUEST_BODY_ERROR_CODE,
			Message: constant.ERR_INVALID_REQUEST_BODY_MESSAGE,
		})
	}

	var validationErr *model.ValidationError

	response, err := controller.UserUsecase.Register(ctx, payload)
	if err != nil {
		if errors.As(err, &validationErr) {
			return util.SendErrorResponse(ctx, err)
		}

		return util.SendErrorResponseInternalServer(ctx, controller.Log, err)
	}

	return util.SendSuccessResponseWithData(ctx, response)
}

func (controller UserController) GetUserInfo(ctx *fiber.Ctx) error {
	userId, err := ctx.ParamsInt("userId")
	if err != nil {
		return util.SendErrorResponseInternalServer(ctx, controller.Log, err)
	}

	var validationErr *model.ValidationError

	response, err := controller.UserUsecase.GetUserInfo(ctx, userId)
	if err != nil {
		if errors.As(err, &validationErr) {
			return util.SendErrorResponseNotFound(ctx, err)
		}

		return util.SendErrorResponseInternalServer(ctx, controller.Log, err)
	}

	return util.SendSuccessResponseWithData(ctx, response)
}
