package handlers

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	handlersMocks "github.com/DyadyaRodya/go-shortener/internal/handlers/mocks"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
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

	delChan := make(chan *usecasesdto.DeleteUserShortURLsRequest, 1024)

	h.handlers = NewHandlers(h.usecases, h.config, delChan)
}

func TestRunHandlersSuite(t *testing.T) {
	suite.Run(t, new(handlersSuite))
}

func (h *handlersSuite) TestNewHandlers() {
	delChan := make(chan *usecasesdto.DeleteUserShortURLsRequest, 1024)

	expected := &Handlers{h.usecases, h.config, delChan}
	if got := NewHandlers(h.usecases, h.config, delChan); !reflect.DeepEqual(got, expected) {
		h.Errorf(errors.New("NewHandlers error"), "NewHandlers() = %v, want %v", got, expected)
	}
}
