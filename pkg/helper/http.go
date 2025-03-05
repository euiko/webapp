package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/elnormous/contenttype"
	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/validator"
	"github.com/ggicci/httpin"
)

type (
	jsonInput struct {
		Payload interface{} `in:"body=json"`
	}

	xmlInput struct {
		Payload interface{} `in:"body=xml"`
	}

	writeResponseConfig struct {
		// to override status code
		status int
	}

	WriteResponseOption func(config *writeResponseConfig)

	FieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	ErrorResponse struct {
		Error       string                `json:"error"`
		FieldErrors map[string]FieldError `json:"field_errors,omitempty"`
	}
)

var (
	ErrInvalidRequest = errors.New("invalid request")

	formMediaTypes = []contenttype.MediaType{
		contenttype.NewMediaType("application/x-www-form-urlencoded"),
		contenttype.NewMediaType("multipart/form-data"),
	}

	jsonMediaTypes = []contenttype.MediaType{
		contenttype.NewMediaType("application/json"),
	}

	xmlMediaTypes = []contenttype.MediaType{
		contenttype.NewMediaType("application/xml"),
		contenttype.NewMediaType("text/xml"),
	}
)

func DecodeRequest(r *http.Request, target interface{}) error {
	var (
		contentType, err = contenttype.GetMediaType(r)
	)

	if err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if contentType.MatchesAny(formMediaTypes...) {
		err = httpin.DecodeTo(r, target)
	} else if contentType.MatchesAny(jsonMediaTypes...) {
		input := jsonInput{target}
		err = httpin.DecodeTo(r, &input)
	} else if contentType.MatchesAny(xmlMediaTypes...) {
		input := xmlInput{target}
		err = httpin.DecodeTo(r, &input)
	}

	if err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validator.Validate(target); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	return nil
}

// ResponseWithStatus overrides the default detected status code for the response
func ResponseWithStatus(status int) WriteResponseOption {
	return func(config *writeResponseConfig) {
		config.status = status
	}
}

func WriteResponse(w http.ResponseWriter, data interface{}, opts ...WriteResponseOption) error {
	var (
		err    error
		config = writeResponseConfig{
			status: 0,
		}
	)

	for _, opt := range opts {
		opt(&config)
	}

	// check for errors
	if err, ok := data.(error); ok {
		errStatus := http.StatusInternalServerError
		body := ErrorResponse{
			Error:       err.Error(),
			FieldErrors: make(map[string]FieldError),
		}
		if errors.Is(err, ErrInvalidRequest) {
			errStatus = http.StatusBadRequest
			body.Error = ErrInvalidRequest.Error()
		}

		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			for _, field := range *validationErrors {
				body.FieldErrors[field.Field()] = FieldError{
					Field:   field.Field(),
					Message: field.Error(),
					Error:   field.Tag(),
				}
			}
		}

		if config.status > 0 {
			errStatus = config.status
		}

		return writeJSON(w, body, errStatus)
	}

	val := reflect.ValueOf(data)
	// dereference pointers
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct, reflect.Map, reflect.Array:
		err = writeJSON(w, data, config.status)
	default:
		err = fmt.Errorf("writing a %s is not supported yet", val.Kind())
		log.Error(err.Error())
	}

	return err
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) error {
	// set default status code
	if status <= 0 {
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
