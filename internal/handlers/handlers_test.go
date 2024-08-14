package handlers

import (
	"errors"
	handlersMocks "github.com/DyadyaRodya/go-shortener/internal/handlers/mocks"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type handlersSuite struct {
	suite.Suite

	usecases *handlersMocks.Usecases
	config   *Config

	handlers *Handlers
}

func (h *handlersSuite) SetupTest() {
	t := h.T()

	h.usecases = handlersMocks.NewUsecases(t)
	h.config = &Config{BaseShortURL: "http://test.example.com"}

	h.handlers = NewHandlers(h.usecases, h.config)
}

func TestRunHandlersSuite(t *testing.T) {
	suite.Run(t, new(handlersSuite))
}

func (h *handlersSuite) TestNewHandlers() {
	expected := &Handlers{h.usecases, h.config}
	if got := NewHandlers(h.usecases, h.config); !reflect.DeepEqual(got, expected) {
		h.Errorf(errors.New("NewHandlers error"), "NewHandlers() = %v, want %v", got, expected)
	}
}
