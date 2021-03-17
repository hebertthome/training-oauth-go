package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/hebertthome/traning-oauth-go/context"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
)

type AppHandler struct {
	Auth bool
	C    *context.AppContext
	H    func(*context.AppContext, http.ResponseWriter, *http.Request) (int, error)
}

// Our ServeHTTP method is mostly the same, and also has the ability to
// access our *appContext's fields (templates, loggers, etc.) as well.
func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ah.Auth {
		if err := ah.authorize(w, r); err != nil {
			ah.C.Logger.Error("Authorize",
				logger.String("Authorize", err.Error()),
			)
			return
		}
	}
	status, err := ah.H(ah.C, w, r)
	if err != nil {
		ah.C.Logger.Info("HTTP",
			logger.Int("status", status),
			logger.Struct("error", err),
		)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, err.Error(), status)
		default:
			http.Error(w, err.Error(), status)
		}
	}
}

func (ah AppHandler) authorize(w http.ResponseWriter, r *http.Request) error {
	// Read attribute '(Bearer) Authorization' on header
	a := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(a) != 2 || a[0] != "Bearer" {
		ah.C.Logger.Info("Authenticate",
			logger.String("(Bearer) Authorization", "Not found"),
		)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("(Bearer) Authorization not found")
	}
	auth := a[1]
	ah.C.Logger.Debug("Authorize",
		logger.String("Cache(key)", auth),
	)
	// Find token authorization on Cache
	result, err := ah.C.Cache.Get(auth)
	if err != nil {
		ah.C.Logger.Error("Authorize",
			logger.Struct("Cache(key) error", err),
		)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return errors.New("Authorization fail!")
	}

	result.Count = result.Count + 1
	ah.C.Cache.Set(auth, *result)

	ah.C.Logger.Debug("Authorize",
		logger.Struct("Cache(Value)", result),
	)
	return nil
}
