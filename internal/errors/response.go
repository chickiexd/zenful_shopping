package errors

import (
	"net/http"

	"github.com/chickiexd/zenful_shopping/internal/logger"
	"github.com/chickiexd/zenful_shopping/utils"
)

// InternalServerError logs and returns a 500 error
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
    logger.Log.Errorw("internal server error",
        "method", r.Method,
        "path", r.URL.Path,
        "error", err.Error(),
    )
    utils.WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

// BadRequest logs and returns a 400 error
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
    logger.Log.Warnw("bad request",
        "method", r.Method,
        "path", r.URL.Path,
        "error", err.Error(),
    )
    utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
}

// NotFound logs and returns a 404 error
func NotFound(w http.ResponseWriter, r *http.Request) {
    logger.Log.Warnw("not found",
        "method", r.Method,
        "path", r.URL.Path,
    )
    utils.WriteJSONError(w, http.StatusNotFound, "resource not found")
}

// Unauthorized logs and returns a 401 error
func Unauthorized(w http.ResponseWriter, r *http.Request) {
    logger.Log.Warnw("unauthorized access attempt",
        "method", r.Method,
        "path", r.URL.Path,
    )
    utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// UnauthorizedBasic logs and returns a 401 with WWW-Authenticate header
func UnauthorizedBasic(w http.ResponseWriter, r *http.Request) {
    logger.Log.Warnw("unauthorized basic auth attempt",
        "method", r.Method,
        "path", r.URL.Path,
    )
    w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
    utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// Forbidden logs and returns a 403 error
func Forbidden(w http.ResponseWriter, r *http.Request) {
    logger.Log.Warnw("forbidden access attempt",
        "method", r.Method,
        "path", r.URL.Path,
    )
    utils.WriteJSONError(w, http.StatusForbidden, "forbidden")
}

// Conflict logs and returns a 409 error
func Conflict(w http.ResponseWriter, r *http.Request, err error) {
    logger.Log.Warnw("resource conflict",
        "method", r.Method,
        "path", r.URL.Path,
        "error", err.Error(),
    )
    utils.WriteJSONError(w, http.StatusConflict, err.Error())
}

// RateLimitExceeded logs and returns a 429 error
func RateLimitExceeded(w http.ResponseWriter, r *http.Request, retryAfter string) {
    logger.Log.Warnw("rate limit exceeded",
        "method", r.Method,
        "path", r.URL.Path,
        "retry_after", retryAfter,
    )
    w.Header().Set("Retry-After", retryAfter)
    utils.WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
