package response

import (
	"blog/pkg/exception"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FailedResponse(ctx *gin.Context, err error, opts ...Option) {
	var (
		errCode  int
		httpCode int
		ns       string
		reason   string
		data     interface{}
		meta     interface{}
	)

	switch t := err.(type) {
	case exception.APIException:
		errCode = t.ErrorCode()
		reason = t.GetReason()
		data = t.GetData()
		meta = t.GetMeta()
		httpCode = t.GetHttpCode()
		ns = t.GetNamespace()
	default:
		errCode = exception.UnKnownException
	}

	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}

	resp := Data{
		Code:      &errCode,
		Namespace: ns,
		Reason:    reason,
		Message:   err.Error(),
		Data:      data,
		Meta:      meta,
	}

	for _, opt := range opts {
		opt.Apply(&resp)
	}
	ctx.JSON(http.StatusOK, resp)
}

func SuccessResponse(ctx *gin.Context, data interface{}, opts ...Option) {
	c := 0
	resp := Data{
		Code:    &c,
		Message: "",
		Data:    data,
	}

	for _, opt := range opts {
		opt.Apply(&resp)
	}
	ctx.JSON(http.StatusOK, resp)
}
