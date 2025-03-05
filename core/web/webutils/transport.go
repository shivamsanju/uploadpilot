package webutils

import (
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
	"github.com/uploadpilot/core/pkg/validator"
)

var validatorSingleton *validator.Validator
var once sync.Once

func NewTransportValidator() *validator.Validator {
	once.Do(func() {
		validatorSingleton = validator.NewValidator()
	})
	return validatorSingleton
}

func CreateJSONHandler[Params any, Query any, Body any, Result any](
	h func(r *http.Request, params Params, query Query, body Body) (Result, int, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validator := NewTransportValidator()

		// Decode params
		var params Params
		if shouldDecode(params) {
			paramsMap := make(map[string]string)
			routeContext := chi.RouteContext(r.Context())
			for i, key := range routeContext.URLParams.Keys {
				paramsMap[key] = routeContext.URLParams.Values[i]
			}
			if err := mapstructure.Decode(paramsMap, &params); err != nil {
				HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			if err := validator.ValidateStruct(params); err != nil {
				HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

		// Decode query
		var query Query
		if shouldDecode(query) {
			queryMap := make(map[string]string)
			for k, v := range r.URL.Query() {
				queryMap[k] = strings.Join(v, ",")
			}
			if err := mapstructure.Decode(queryMap, &query); err != nil {
				HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			if err := validator.ValidateStruct(query); err != nil {
				HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

		// Decode body
		var body Body
		if shouldDecode(body) {
			if err := render.DecodeJSON(r.Body, &body); err != nil && r.ContentLength > 0 {
				HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			if err := validator.ValidateStruct(body); err != nil {
				HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

		// Call the handler function
		result, status, err := h(r, params, query, body)
		if err != nil {
			HandleHttpError(w, r, status, err)
			return
		}

		render.JSON(w, r, result)
	}
}

func shouldDecode(v any) bool {
	if v == nil {
		return false
	}
	typeName := reflect.TypeOf(v).String()
	return typeName != "interface {}" && typeName != "struct {}" && typeName != "<nil>"
}
