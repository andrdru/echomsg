package echomsg

import (
	"encoding/json"
	"net/http"
	"strings"

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

// SetError set error with list and map
func (m *messageJson) SetError(options ...Option) {
	var args = &Options{}

	var opt Option
	for _, opt = range options {
		opt(args)
	}

	if m.ErrorCode == 0 ||
		m.ErrorCode == http.StatusOK ||
		m.ErrorCode < http.StatusInternalServerError && args.code >= http.StatusInternalServerError {
		m.ErrorCode = args.code
	}

	if args.field == "" && args.message == "" {
		return
	}

	var parts = make([]string, 0, 2)
	if args.field != "" {
		parts = append(parts, args.field)
	}
	if args.message != "" {
		parts = append(parts, args.message)
	}

	m.ErrorMessages = append(m.ErrorMessages, strings.Join(parts, ": "))

	if len(parts) == 2 {
		if m.ErrorMaps == nil {
			m.ErrorMaps = make(map[string][]string)
		}

		m.ErrorMaps[args.field] = append(m.ErrorMaps[args.field], args.message)
	}
}

func (m *messageJson) Raw() interface{} {
	return m
}
