package responder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Responder interface {
	ErrorBedRequest(ctx *gin.Context, err error)
	ErrorInternal(ctx *gin.Context, err error)
	ErrorNotFound(ctx *gin.Context, err error)
	SuccessJSON(ctx *gin.Context, responseData interface{})
}

type respond struct {
	log *zap.Logger
}

func NewRespond(logger *zap.Logger) Responder {
	return &respond{log: logger}
}

func (r *respond) ErrorBedRequest(ctx *gin.Context, err error) {
	r.log.Info("http response bad request", zap.Error(err))
	ctx.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}

func (r *respond) ErrorInternal(ctx *gin.Context, err error) {
	r.log.Error("http ErrorResponse internal error", zap.Error(err))
	ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

func (r *respond) ErrorNotFound(ctx *gin.Context, err error) {
	r.log.Error("http ErrorResponse not found", zap.Error(err))
	ctx.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: err.Error(),
	})
}

func (r *respond) SuccessJSON(ctx *gin.Context, responseData interface{}) {
	r.log.Info("Successful operation, sending back response data")
	ctx.JSON(http.StatusOK, responseData)
}
