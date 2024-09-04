package responder

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"user-service/internal/models"
)

type ResponderInterface interface {
	SuccessJSON(ctx *gin.Context, responseData interface{})
	BadRequest(ctx *gin.Context, err error)
	ErrorInternal(ctx *gin.Context, err error)
	NotAuthorized(ctx *gin.Context, err error)
}

type ResponderObj struct {
	logger *zap.SugaredLogger
}

func NewResponderObj() *ResponderObj {
	return &ResponderObj{logger: zap.NewExample().Sugar()}
}

func (r *ResponderObj) SuccessJSON(ctx *gin.Context, responseData interface{}) {
	r.logger.Info("Successful operation, sending back response data")
	ctx.JSON(http.StatusOK, responseData)
}

func (r *ResponderObj) BadRequest(ctx *gin.Context, err error) {
	r.logger.Info("Bad request", zap.Error(err))
	ctx.JSON(http.StatusBadRequest, models.Response{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}

func (r *ResponderObj) ErrorInternal(ctx *gin.Context, err error) {
	r.logger.Info("Internal Error", zap.Error(err))
	ctx.JSON(http.StatusInternalServerError, models.Response{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

func (r *ResponderObj) NotAuthorized(ctx *gin.Context, err error) {
	r.logger.Info("Not authorized", zap.Error(err))
	ctx.JSON(http.StatusForbidden, models.Response{
		Code:    http.StatusForbidden,
		Message: err.Error(),
	})
}
