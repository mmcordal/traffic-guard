package errorsx

import (
	"errors"
	"net/http"
	"strings"
	"traffic-guarder/internal/viewmodel"

	"github.com/gofiber/fiber/v2"
)

type APIError interface {
	Error() string
	x()
}

type ErrorType int

const (
	ErrorTypeBadRequest ErrorType = iota + 1
	ErrorTypeUnauthorized
	ErrorTypeNotFound
	ErrorTypeDatabase
	ErrorTypeInternal
)

type ErrorCode int

const (
	ErrorCodeBadRequest   ErrorCode = 400
	ErrorCodeUnauthorized ErrorCode = 401
	ErrorCodeNotFound     ErrorCode = 404
	ErrorCodeDataBase     ErrorCode = 500
	ErrorCodeInternal     ErrorCode = 500
)

type CustomError struct {
	message   error
	ErrorType ErrorType
	ErrorCode ErrorCode
}

func (CustomError) x() {}

func (error CustomError) Error() string {
	return error.message.Error()
}

func ValidationError(errs []error) APIError {
	var sb strings.Builder
	for _, err := range errs {
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return CustomError{message: errors.New(sb.String()), ErrorType: ErrorTypeBadRequest, ErrorCode: ErrorCodeBadRequest}
}

func UnauthorizedError(err error) APIError {
	return CustomError{err, ErrorTypeUnauthorized, ErrorCodeUnauthorized}
}

func NotFoundError(err error) APIError {
	return CustomError{err, ErrorTypeNotFound, ErrorCodeNotFound}
}

func DatabaseError(err error) APIError {
	return CustomError{err, ErrorTypeDatabase, ErrorCodeDataBase}
}

func InternalError(err error) APIError {
	return CustomError{err, ErrorTypeInternal, ErrorCodeInternal}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	respModel := viewmodel.ErrorResponse{
		Success:      false,
		ErrorMessage: "error",
	}

	var customErr CustomError
	if !errors.As(err, &customErr) {
		respModel.ErrorMessage = "Bilinmeyen sunucu hatası: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(respModel)
	}

	switch customErr.ErrorType {
	case ErrorTypeBadRequest:
		respModel.ErrorMessage = "Hatalı istek: " + err.Error()
		return c.Status(http.StatusBadRequest).JSON(respModel)
	case ErrorTypeUnauthorized:
		respModel.ErrorMessage = "Yetkisiz erişim " + err.Error()
		return c.Status(http.StatusUnauthorized).JSON(respModel)
	case ErrorTypeNotFound:
		respModel.ErrorMessage = "Kayıt bulunamadı: " + err.Error()
		return c.Status(http.StatusNotFound).JSON(respModel)
	case ErrorTypeInternal:
		respModel.ErrorMessage = "Sunucu iç hatası: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(respModel)
	case ErrorTypeDatabase:
		respModel.ErrorMessage = "Sunucu hatası: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(respModel)
	}
	return c.Status(http.StatusInternalServerError).JSON(respModel)
}
