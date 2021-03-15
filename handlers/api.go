package handlers

import (
	"fmt"
	"net/http"

	"bitbucket.org/hebertthome/traning-oauth-go/context"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
)

func API(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	// Build Response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("{\"result\": \"%s\"}", "Success")))
	ctx.Logger.Info("API",
		logger.String("Request", "Success"),
	)
	return http.StatusOK, nil
}
