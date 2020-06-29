package echomsg

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	MessageJson interface {
		Message
		Error(code int, message string) *messageJson
		SetError(code int, message string)
	}

	messageJson struct {
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
var _ MessageJson = &messageJson{}

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
	for _, el = range m.ErrorMessages {
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

func (m *messageJson) Error(code int, message string) *messageJson {
	var c = m.cloneErr()

	c.SetError(code, message)
	return c
}

func (m *messageJson) SetError(code int, message string) {
	if m.ErrorCode == 0 ||
		m.ErrorCode == http.StatusOK ||
		m.ErrorCode < http.StatusInternalServerError && code >= http.StatusInternalServerError {
		m.ErrorCode = code
	}

	m.ErrorMessages = append(m.ErrorMessages, message)
}
