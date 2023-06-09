package api

import (
	"errors"
	"fmt"
	"github.com/a-novel/agora-backend/environment"
	"github.com/a-novel/agora-backend/framework/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"
)

const (
	ContextData = "agora_ctx"
)

type ErrWithStatus struct {
	Err  error
	Code int
}

func ErrToStatus(err error, st []ErrWithStatus, masks map[error]int) int {
	for _, v := range st {
		if errors.Is(err, v.Err) {
			// Use the mask code, if provided.
			if masks != nil {
				if code, ok := masks[v.Err]; ok {
					return code
				}
			}

			return v.Code
		}
	}

	return http.StatusInternalServerError
}

type CallbackResponse struct {
	Body                 interface{}
	CTXData              interface{}
	MaskErrorsWithStatus map[error]int
	Deferred             environment.Deferred
}

type Callback[Body any, Provider any] func(c *gin.Context, token string, body Body, provider Provider) (CallbackResponse, error)

func WithContext[Body any, Provider any](callback Callback[Body, Provider], provider Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body Body

		token := c.GetHeader("Authorization")

		// Parse body if one is expected.
		if reflect.TypeOf(body) != nil {
			// Ignore EOF errors, because it means there is no input stream and this is normal.
			if err := c.ShouldBindUri(&body); err != nil && !errors.Is(err, io.EOF) {
				_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
				return
			}
			if err := c.ShouldBindJSON(&body); err != nil && !errors.Is(err, io.EOF) {
				_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
				return
			}
		}

		resp, err := callback(c, token, body, provider)
		if err != nil {
			status := ErrToStatus(err, []ErrWithStatus{
				{Err: validation.ErrUniqConstraintViolation, Code: http.StatusConflict},
				{Err: validation.ErrConstraintViolation, Code: http.StatusUnprocessableEntity},
				{Err: validation.ErrTimeout, Code: http.StatusRequestTimeout},
				{Err: validation.ErrInvalidEntity, Code: http.StatusUnprocessableEntity},
				{Err: validation.ErrInvalidCredentials, Code: http.StatusForbidden},
				{Err: validation.ErrUnauthorized, Code: http.StatusUnauthorized},
				{Err: validation.ErrValidated, Code: http.StatusGone},
				{Err: validation.ErrNotFound, Code: http.StatusNotFound},
			}, resp.MaskErrorsWithStatus)
			_ = c.AbortWithError(status, err)
		}

		if resp.CTXData != nil {
			c.Set(ContextData, resp.CTXData)
		}

		// A deferred function can still be returned, even if the main handler throws an error.
		if resp.Deferred != nil {
			// TODO: implement retry.
			if err = resp.Deferred(); err != nil {
				_ = c.Error(fmt.Errorf("error while executing deferred function: %w", err))
			}
		}

		if err == nil {
			if resp.Body == nil {
				c.Status(http.StatusNoContent)
			} else {
				c.JSON(http.StatusOK, resp.Body)
			}
		}

		c.Next()
	}
}

type HandlerMap map[string]gin.HandlerFunc

type Config map[string]HandlerMap

func LoadAPI(r gin.IRouter, basePath string, routes Config) {
	for route, cfg := range routes {
		for method, handler := range cfg {
			r.Handle(method, path.Join(basePath, route), handler)
		}
	}
}

func Logger(logger zerolog.Logger, projectID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		status := c.Writer.Status()
		errs := c.Errors.Errors()

		logLevel := zerolog.TraceLevel
		severity := "INFO" // For GCP.
		if status > 499 {
			logLevel = zerolog.ErrorLevel
			severity = "ERROR"
		} else if status > 399 || len(errs) > 0 {
			logLevel = zerolog.WarnLevel
			severity = "WARNING"
		}

		parserQuery := zerolog.Dict()
		for k, v := range c.Request.URL.Query() {
			parserQuery.Strs(k, v)
		}

		// Allow logs to be grouped in log explorer.
		// https://cloud.google.com/run/docs/logging#run_manual_logging-go
		var trace string
		if projectID != "" {
			traceHeader := c.GetHeader("X-Cloud-Trace-Context")
			traceParts := strings.Split(traceHeader, "/")
			if len(traceParts) > 0 && len(traceParts[0]) > 0 {
				trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
			}
		}

		ll := logger.WithLevel(logLevel).
			Dict(
				"httpRequest", zerolog.Dict().
					Str("requestMethod", c.Request.Method).
					Str("requestUrl", c.FullPath()).
					Int("status", status).
					Str("userAgent", c.Request.UserAgent()).
					Str("remoteIp", c.ClientIP()).
					Str("protocol", c.Request.Proto).
					Str("latency", end.Sub(start).String()),
			).
			Time("start", start).
			Dur("postProcessingLatency", time.Now().Sub(end)).
			Int64("contentLength", c.Request.ContentLength).
			Str("ip", c.ClientIP()).
			Str("contentType", c.ContentType()).
			Str("auth", c.GetHeader("Authorization")).
			Strs("errors", errs).
			Str("severity", severity).
			Interface("context", c.GetStringMap(ContextData))

		if len(trace) > 0 {
			ll = ll.Str("logging.googleapis.com/trace", trace)
		}

		ll.Msg(c.Request.URL.String())
	}
}
