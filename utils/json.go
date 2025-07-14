package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func ParseRequestBody(c echo.Context, data interface{}) error {

	// body, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	logrus.Errorf("ErrorParsing request body: %v", err)
	// 	// c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
	// 	return err
	// }

	// if len(body) > 0 {
	err := c.Bind(data)
	if err != nil {
		logrus.Errorf("ErrorParsing request body: %v", err)
		return err
	}
	// }

	return err
}

func Decode(data string, destination interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.DisallowUnknownFields() // <-- MAGIC DISINI 🔥

	if err := decoder.Decode(&destination); err != nil {
		return err
	}

	return nil
}

func WriteCustomResponse(c echo.Context, data interface{}, err error) {
	var code int = 200
	if err != nil {
		if _, ok := err.(*BadRequestError); ok {
			code = 400
		} else if _, ok := err.(*NotFoundError); ok {
			code = 404
		} else if _, ok := err.(*UnprocessableContentError); ok {
			code = 422
		} else if _, ok := err.(*ConflictError); ok {
			code = 409
		} else if _, ok := err.(*UnauthorizedError); ok {
			code = 401
		} else {
			code = 500
		}
	}

	c.JSON(code, data)
}

func ValidateFilter(source string, destination interface{}) error {
	return json.Unmarshal([]byte(source), &destination)
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
	response := struct {
		Status     string `json:"status"`
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}{
		Status:     "error",
		StatusCode: code,
		Message:    msg,
	}

	// Status:     "error",
	// 	StatusCode: code,
	// 	Message:    msg,

	// Send the JSON response
	if !c.Response().Committed {
		c.JSON(code, response)
	}
}
