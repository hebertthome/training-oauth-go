package handlers

import (
	"errors"
	"net/http"

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
	// Check header attribute named 'Authorize'
	auth := r.Header.Get("Authorize")
	if len(auth) < 1 {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return errors.New("Header(Authorize) not found!")
	}
	ah.C.Logger.Debug("Authorize",
		logger.String("Cache(key)", auth),
	)
	// Find token authorization on Cache
	result, err := ah.C.Cache.Get(auth)
	if err != nil {
		ah.C.Logger.Error("Authorize",
			logger.Struct("Cache(key)", err),
		)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return errors.New("Authorization fail!")
	}

	result.Count = result.Count + 1
	ah.C.Cache.Set(auth, *result)

	ah.C.Logger.Debug("Authorize",
		logger.Struct("Cache(Value)", result),
	)
	return nil
}
