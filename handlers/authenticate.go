package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"bitbucket.org/hebertthome/traning-oauth-go/config"
	"bitbucket.org/hebertthome/traning-oauth-go/context"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
	"bitbucket.org/hebertthome/traning-oauth-go/redis"
)

func Authenticate(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	// Read attribute '(Basic) Authorization' on header
	a := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(a) != 2 || a[0] != "Basic" {
		ctx.Logger.Info("Authenticate",
			logger.String("(Basic) Authentication", "Not found"),
		)
		return http.StatusUnauthorized, fmt.Errorf("(Basic) Authentication not found")
	}
	payload, _ := base64.StdEncoding.DecodeString(a[1])
	pair := strings.SplitN(string(payload), ":", 2)
	// Load static credentials :)
	auth := config.Get().Auth
	// Compare credetionais
	if pair[0] != auth.ClientID || pair[1] != auth.ClientSecret {
		ctx.Logger.Info("Authenticate",
			logger.String("credentials", "Wrong!"),
		)
		return http.StatusUnauthorized, fmt.Errorf("wrong credentials :(")
	}
	out := uuid.New()
	// Write auth token on cache
	ctx.Cache.Set(out.String(), redis.LcCached{ID: pair[0], Count: 0})
	ctx.Logger.Info("Authenticate",
		logger.String("Token", out.String()),
	)
	// Build response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("{\"access_token\": \"%s\"}", out.String())))
	ctx.Logger.Info("Authenticate",
		logger.String("token", "Success generated"),
	)
	return http.StatusOK, nil
}
