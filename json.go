package echomsg

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	messageJson struct {
		Message
		Data          interface{}
		ErrorCode     int
		ErrorMessages []string
	}

	messageJsonError struct {
		ErrorCode     int      `json:"code,omitempty"`
		ErrorMessages []string `json:"messages,omitempty"`
	}
)

var _ json.Marshaler = &messageJson{}
var _ Message = &messageJson{}

func NewJson() *messageJson {
	return &messageJson{
		ErrorCode: http.StatusOK,
	}
}

func (m *messageJson) cloneErr() *messageJson {
	var c = &messageJson{
		Data:      m.Data,
		ErrorCode: m.ErrorCode,
	}

	var el string
	for _, el = range c.ErrorMessages {
		c.ErrorMessages = append(c.ErrorMessages, el)
	}

	return c
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
	})
}

func (m *messageJson) Return(c echo.Context) error {
	return c.JSON(m.ErrorCode, m)
}

func (m *messageJson) Error(code int, message string) Message {
	var c = m.cloneErr()
	if c.ErrorCode == 0 ||
		c.ErrorCode == http.StatusOK ||
		c.ErrorCode < http.StatusInternalServerError && code >= http.StatusInternalServerError {
		c.ErrorCode = code
	}

	c.ErrorMessages = append(c.ErrorMessages, message)

	return c
}
