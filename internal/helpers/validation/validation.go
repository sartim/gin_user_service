package validation

import (
	"errors"
	"gin-shop-api/internal/helpers/types"
	"strings"

	"github.com/go-playground/validator/v10"
)

func MsgForTag(tag string, fieldType string) string {
	switch tag {
	case "required":
		switch fieldType {
		case "body":
			return "This field is required"
		case "header":
			return "Missing header"
		}
	}
	return ""
}

func ValidateSchema(err error, fieldType string) []types.Dict {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		body := []types.Dict{}
		for _, fe := range ve {
			errors := types.Dict{}
			if fe.Field() == "ContentType" {
				errors["Content-Type"] = MsgForTag(fe.Tag(), fieldType)
			} else if fe.Field() == "Authorization" {
				errors[fe.Field()] = MsgForTag(fe.Tag(), fieldType)
			} else {
				errors[strings.ToLower(fe.Field())] = MsgForTag(fe.Tag(), fieldType)
			}
			body = append(body, errors)
		}
		return body
	}
	return nil
}
