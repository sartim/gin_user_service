package helpers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

func ValidateSchema(ctx *gin.Context, err error, fieldType string) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		body := []Dict{}
		for _, fe := range ve {
			errors := Dict{}
			if fe.Field() == "ContentType" {
				errors["Content-Type"] = MsgForTag(fe.Tag(), fieldType)
			} else if fe.Field() == "Authorization" {
				errors[fe.Field()] = MsgForTag(fe.Tag(), fieldType)
			} else {
				errors[strings.ToLower(fe.Field())] = MsgForTag(fe.Tag(), fieldType)
			}

			body = append(body, errors)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": body})
	}
}
