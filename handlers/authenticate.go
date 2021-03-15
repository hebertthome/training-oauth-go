package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"bitbucket.org/hebertthome/traning-oauth-go/config"
	"bitbucket.org/hebertthome/traning-oauth-go/context"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
	"bitbucket.org/hebertthome/traning-oauth-go/redis"
)

func Authenticate(ctx *context.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	// Read attribute 'client-id' on query string
	id, ok := r.URL.Query()["client-id"]
	if !ok || len(id[0]) < 1 {
		ctx.Logger.Info("Authenticate",
			logger.String("client-id", "Not found"),
		)
		return http.StatusUnauthorized, fmt.Errorf("client-id not found")
	}
	// Read attribute 'client-secret' on query string
	secret, ok := r.URL.Query()["client-secret"]
	if !ok || len(secret[0]) < 1 {
		ctx.Logger.Info("Authenticate",
			logger.String("client-secret", "Not found"),
		)
		return http.StatusUnauthorized, fmt.Errorf("client-secret not found")
	}
	// Load static credentials :)
	auth := config.Get().Auth
	// Compare credetionais
	if id[0] != auth.ClientID || secret[0] != auth.ClientSecret {
		ctx.Logger.Info("Authenticate",
			logger.String("credentials", "Wrong!"),
		)
		return http.StatusUnauthorized, fmt.Errorf("wrong credentials :(")
	}
	out := uuid.New()
	// Write auth token on cache
	ctx.Cache.Set(out.String(), redis.LcCached{ID: id[0], Count: 0})
	// Build response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("{\"token\": \"%s\"}", out.String())))
	ctx.Logger.Info("Authenticate",
		logger.String("token", "Success generated"),
	)
	return http.StatusOK, nil
}
