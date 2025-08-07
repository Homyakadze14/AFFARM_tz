package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ParseErr(errs error) (status int, err error) {
	if errs == nil {
		return http.StatusOK, nil
	}

	newErrMes := ""
	status = 0

	var ve validator.ValidationErrors
	if errors.As(errs, &ve) {
		status = http.StatusBadRequest
		for _, v := range ve {
			if v.Tag() == "required" {
				newErrMes += fmt.Sprintf("Field %s must be provided;", v.Field())
			}
			if v.Tag() == "email" {
				newErrMes += fmt.Sprintf("Field %s must contains email;", v.Field())
			}
			if v.Tag() == "min" {
				newErrMes += fmt.Sprintf("Minimal lenght for field %s is %v;", v.Field(), v.Param())
			}
			if v.Tag() == "max" {
				newErrMes += fmt.Sprintf("Maximum lenght for field %s is %v;", v.Field(), v.Param())
			}
		}
	}

	switch {
	case errors.Is(errs, ErrCryptocurrencyAlreadyExists):
		status = http.StatusBadRequest
		newErrMes += fmt.Sprintf("%v;", ErrCryptocurrencyAlreadyExists)
	case errors.Is(errs, ErrCryptocurrencyNotFound):
		status = http.StatusNotFound
		newErrMes += fmt.Sprintf("%v;", ErrCryptocurrencyNotFound)
	case errors.Is(errs, ErrHistoryNotFound):
		status = http.StatusNotFound
		newErrMes += fmt.Sprintf("%v;", ErrHistoryNotFound)
	case errors.Is(errs, ErrTrackingAlreadyExists):
		status = http.StatusBadRequest
		newErrMes += fmt.Sprintf("%v;", ErrTrackingAlreadyExists)
	case errors.Is(errs, ErrTrackingNotFound):
		status = http.StatusNotFound
		newErrMes += fmt.Sprintf("%v;", ErrTrackingNotFound)
	case errors.Is(errs, ErrSymbolNotFound):
		status = http.StatusNotFound
		newErrMes += fmt.Sprintf("%v;", ErrSymbolNotFound)
	case newErrMes == "":
		status = http.StatusInternalServerError
		newErrMes = "Internal server error"
	}

	return status, errors.New(newErrMes)
}
