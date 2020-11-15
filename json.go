package echomsg

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:generate easyjson

type (
	messageJson struct {
		Data          interface{}
		ErrorCode     int
		ErrorMessages []string
		ErrorMaps     map[string][]string
	}

	//easyjson:json
	messageJsonError struct {
		ErrorCode     int                 `json:"code"`
		ErrorMessages []string            `json:"messages"`
		ErrorMaps     map[string][]string `json:"maps"`
	}
)

var _ json.Marshaler = &messageJson{}
var _ Message = &messageJson{}

func NewJson() *messageJson {
	return &messageJson{
		ErrorCode: http.StatusOK,
	}
}

func (m *messageJson) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	if m.ErrorCode == http.StatusOK {
		if m.Data == nil {
			m.Data = struct{}{}
		}
		return json.Marshal(m.Data)
	}

	if m.ErrorCode == 0 {
		m.ErrorCode = http.StatusInternalServerError
	}

	return json.Marshal(messageJsonError{
		ErrorCode:     m.ErrorCode,
		ErrorMessages: m.ErrorMessages,
		ErrorMaps:     m.ErrorMaps,
	})
}

func (m *messageJson) Return(c echo.Context) error {
	return c.JSON(m.ErrorCode, m)
}

// SetError set error message
// deprecated use SetErrorMap instead
func (m *messageJson) SetError(code int, message string) {
	if m.ErrorCode == 0 ||
		m.ErrorCode == http.StatusOK ||
		m.ErrorCode < http.StatusInternalServerError && code >= http.StatusInternalServerError {
		m.ErrorCode = code
	}

	if message != "" {
		m.ErrorMessages = append(m.ErrorMessages, message)
	}
}

// SetErrorMap set error with list and map
func (m *messageJson) SetErrorMap(code int, message string, field string) {
	if m.ErrorCode == 0 ||
		m.ErrorCode == http.StatusOK ||
		m.ErrorCode < http.StatusInternalServerError && code >= http.StatusInternalServerError {
		m.ErrorCode = code
	}

	if message != "" {
		m.ErrorMessages = append(m.ErrorMessages, message)
	}

	if field != "" {
		if m.ErrorMaps == nil {
			m.ErrorMaps = make(map[string][]string)
		}

		m.ErrorMaps[field] = append(m.ErrorMaps[field], message)
	}
}
