package web

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	vald "github.com/uploadpilot/go-core/common/validator"
	"github.com/uploadpilot/manager/internal/utils"
)

func CreateJSONHandler[Params any, Query any, Body any, Result any](
	h func(r *http.Request, params Params, query Query, body Body) (Result, int, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Decode params
		var params Params
		if shouldDecode(params) {
			paramsMap := make(map[string]string)
			routeContext := chi.RouteContext(r.Context())
			for i, key := range routeContext.URLParams.Keys {
				paramsMap[key] = routeContext.URLParams.Values[i]
			}
			if err := mapstructure.Decode(paramsMap, &params); err != nil {
				utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			if err := validateStruct(params); err != nil {
				utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
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
				utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			if err := validateStruct(query); err != nil {
				utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

		// Decode body
		var body Body
		if shouldDecode(body) {
			if err := render.DecodeJSON(r.Body, &body); err != nil && r.ContentLength > 0 {
				utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			if err := validateStruct(body); err != nil {
				utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

		// Call the handler function
		result, status, err := h(r, params, query, body)
		if err != nil {
			utils.HandleHttpError(w, r, status, err)
			return
		}

		render.JSON(w, r, result)
	}
}

func validateStruct(s any) error {
	validate := vald.NewValidator()
	if err := validate.Struct(s); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, e := range ve {
				errMsgs = append(errMsgs, fmt.Sprintf("%s (%s%s)", e.Field(), e.Tag(), e.Param()))
			}
			return errors.New(strings.Join(errMsgs, "; "))
		}
		return err
	}
	return nil
}

// shouldDecode checks if the type is not an empty struct or interface{}
func shouldDecode(v any) bool {
	if v == nil {
		return false
	}
	typeName := reflect.TypeOf(v).String()
	return typeName != "interface {}" && typeName != "struct {}" && typeName != "<nil>"
}
