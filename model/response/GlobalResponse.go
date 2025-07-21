package response

import (
	"booking/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type GlobalJSONResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type GlobalListDataResponse struct {
	MetadataResponse
	List []any `json:"list"`
}

type GlobalListResponse struct {
	Status     string                 `json:"status"`
	StatusCode string                 `json:"statusCode"`
	Message    string                 `json:"message"`
	Data       GlobalListDataResponse `json:"data,omitempty"`
}

type GlobalSingleResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}
type MetadataResponse struct {
	Page      int `json:"page"`
	Limit     int `json:"perPage"`
	TotalPage int `json:"totalPage"`
	Count     int `json:"totalData"`
}

func StrPtr(s string) *string {
	if s == "" || s == "<nil>" {
		return nil // Mengembalikan nil jika string kosong
	}
	return &s
}

func WriteResponseListJSON(c echo.Context, data GlobalListDataResponse, err error) {
	var code int = 200
	var status string = "OK"
	var message string = "success"

	if err != nil {
		if badRequest, ok := err.(*utils.BadRequestError); ok {
			code = 400
			status = "BadRequest"
			message = badRequest.Message
		} else if notFound, ok := err.(*utils.BadRequestError); ok {
			code = 404
			status = "NotFound"
			message = notFound.Message
		} else {
			code = 500
			status = "ERROR"
			message = err.Error()
		}
	} else {

	}
	response := GlobalListResponse{
		StatusCode: "",
		Status:     status,
		Message:    message,
		Data:       data,
	}

	c.JSON(code, response)
}

func WriteResponseSingleJSON(c echo.Context, data interface{}, err error) {
	var code int = 200
	var status string = "OK"
	var message string = "success"

	if err != nil {
		if badRequest, ok := err.(*utils.BadRequestError); ok {
			code = 400
			status = "BadRequest"
			message = badRequest.Message
		} else if notFound, ok := err.(*utils.NotFoundError); ok {
			code = 404
			status = "NotFound"
			message = notFound.Message
		} else if unprocessable, ok := err.(*utils.UnprocessableContentError); ok {
			code = 422
			status = "UnprocessableEntity"
			message = unprocessable.Message
		} else if conflict, ok := err.(*utils.ConflictError); ok {
			code = 409
			status = "ConflictError"
			message = conflict.Message
		} else if unauthorized, ok := err.(*utils.UnauthorizedError); ok {
			code = 401
			status = "Unauthorized"
			message = unauthorized.Message
		} else if unauthorized, ok := err.(*utils.ForbiddenError); ok {
			code = 403
			status = "Forbidden"
			message = unauthorized.Message
		} else {
			code = 500
			status = "InternalServerError"
			message = "Internal Server Error"
			logrus.Errorf("Internal Server Error : %v", err.Error())
		}
	}
	response := GlobalSingleResponse{
		Status:     status,
		Message:    message,
		StatusCode: strconv.Itoa(code),
		Data:       data,
	}

	c.JSON(code, response)
}
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"

	// If the error is an *echo.HTTPError, use its Code and Message
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		// Ambil message, kalau dia string langsung kasih
		if m, ok := he.Message.(string); ok {
			msg = m
		} else {
			// Kalau bukan string, convert ke string
			msg = http.StatusText(he.Code)
		}
	}

	// Custom JSON structure
	response := GlobalJSONResponse{
		Status:     "error",
		StatusCode: code,
		Message:    msg,
	}

	// Send the JSON response
	if !c.Response().Committed {
		c.JSON(code, response)
	}
}
