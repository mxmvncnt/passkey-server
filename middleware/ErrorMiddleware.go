package middleware

import (
	"errors"
	"net/http"
	"passkey-server/utils/apierror"
	"passkey-server/utils/logger"
	"runtime/debug"

	"github.com/jackc/pgx/v5/pgconn"
)

func ErrorHandler(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			switch err := err.(type) {
			case *apierror.ApiError:
				err.Send(w)
				logger.Errorf("[ErrorMiddleware] ApiError thrown: %v", *err)
			default:
				logger.Errorf("Unexpected exception caught: %v (%s %s)", err, r.Method, r.URL.Path)
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					logger.Errorf("Postgres error: code=%s message=%s detail=%s hint=%s where=%s", pgErr.Code, pgErr.Message, pgErr.Detail, pgErr.Hint, pgErr.Where)
				}
				logger.Errorf("Stack Trace:\n%s\n", debug.Stack())
				apierror.NewApiError(
					http.StatusInternalServerError,
					"UNEXPECTED_ERROR",
					"Please wait and try again",
					"An unexpected internal error has occurred").Send(w)
			}
		}
	}
}
