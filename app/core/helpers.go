package core

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func LogEvent(logLevel string) *log.Logger {
	// Info writes logs in the color blue with "INFO: " as prefix
	var Info = log.New(os.Stdout, "\u001b[34mINFO: \u001B[0m", log.LstdFlags|log.Lshortfile)

	// Warning writes logs in the color yellow with "WARNING: " as prefix
	var Warning = log.New(os.Stdout, "\u001b[33mWARNING: \u001B[0m", log.LstdFlags|log.Lshortfile)

	// Error writes logs in the color red with "ERROR: " as prefix
	var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)

	// Debug writes logs in the color cyan with "DEBUG: " as prefix
	var Debug = log.New(os.Stdout, "\u001b[36mDEBUG: \u001B[0m", log.LstdFlags|log.Lshortfile)

	switch logLevel {
	case "INFO":
		return Info
	case "WARNING":
		return Warning
	case "ERROR":
		return Error
	case "DEBUG":
		return Debug
	}

	return nil
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Log(logLevel string) *log.Logger {
	var logEvent = LogEvent(logLevel)

	return logEvent
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

func GenerateUUID() uuid.UUID {
	return uuid.New()
}
