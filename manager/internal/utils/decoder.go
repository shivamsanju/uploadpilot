package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func CreateJSONHandler[W any](h func(w http.ResponseWriter, r *http.Request) (W, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := h(w, r)
		if err != nil {
			HandleHttpError(w, r, http.StatusBadRequest, err)
			return
		}

		render.JSON(w, r, result)
	}
}

func CreateJSONHandlerWithBody[R any, W any](h func(w http.ResponseWriter, r *http.Request, body *R) (W, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body R

		// Decode JSON request body
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			HandleHttpError(w, r, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnprocessableEntity)))
			return
		}

		// Validate request body
		validate := validator.New()
		if err := validate.Struct(body); err != nil {
			if ve, ok := err.(validator.ValidationErrors); ok {
				errMsgs := ""
				for _, e := range ve {
					errMsgs += fmt.Sprintf("validation failed for field '%s' (%s", e.Field(), e.Tag())
					if e.Param() != "" {
						errMsgs += fmt.Sprintf(" : %s", e.Param())
					}
					errMsgs += "); "
				}
				HandleHttpError(w, r, http.StatusBadRequest, errors.New(errMsgs))
				return
			}
			HandleHttpError(w, r, http.StatusBadRequest, err)
			return
		}

		// Call the handler function
		result, err := h(w, r, &body)
		if err != nil {
			HandleHttpError(w, r, http.StatusBadRequest, err)
			return
		}

		render.JSON(w, r, result)
	}
}
