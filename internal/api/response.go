package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJSON = "application/json"
	MimeTextYAML        = "text/yaml"
)

type (
	responder struct {
		logger *zap.Logger
	}
	apiError struct {
		Error string `json:"errors"`
	}
)

func NewResponder(logger *zap.Logger) *responder {
	return &responder{logger: logger}
}

func (rm *responder) Write(w http.ResponseWriter, code int, mime string, data []byte) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(HeaderContentType, mime)
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		rm.logger.Error("failed to write response", zap.Error(err))
	}
}

func (rm *responder) JSON(w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(HeaderContentType, MimeApplicationJSON)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		rm.logger.Error("failed to response with JSON", zap.Error(err))
	}
}

func (rm *responder) ErrorWithCode(w http.ResponseWriter, err error, code int) {
	if code >= http.StatusInternalServerError {
		rm.logger.Error("got an api errors", zap.Error(err))
		w.WriteHeader(code)
		return
	}

	rm.JSON(w, code, apiError{Error: err.Error()})
}

func (rm *responder) OK(w http.ResponseWriter, data interface{}) {
	rm.JSON(w, http.StatusOK, data)
}

func (rm *responder) Created(w http.ResponseWriter, data interface{}) {
	rm.JSON(w, http.StatusCreated, data)
}

func (rm *responder) NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
