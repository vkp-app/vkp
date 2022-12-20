package promutil

import (
	"context"
	"github.com/djcass44/go-utils/utilities/httputils"
	"net/http"
)

func RespondErr(ctx context.Context, w http.ResponseWriter, err error, code int) {
	httputils.ReturnJSON(ctx, w, code, &Response{
		Status:    "error",
		Data:      struct{}{},
		ErrorType: http.StatusText(code),
		Error:     err.Error(),
		Warnings:  []string{},
	})
}
