package handlers

import (
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	handlersMocks "github.com/DyadyaRodya/go-shortener/internal/handlers/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func (h *handlersSuite) TestGetShortURL() {
	type usecaseResult struct {
		shortURL *entity.ShortURL
		err      error
	}
	type want struct {
		code    int
		headers map[string]string
	}
	tests := []struct {
		name         string
		request      *http.Request
		usecaseParam string
		usecaseRes   *usecaseResult
		want         want
	}{
		{
			name:         "Success",
			request:      httptest.NewRequest(http.MethodGet, "/10abcdef", nil),
			usecaseParam: "10abcdef",
			usecaseRes: &usecaseResult{
				shortURL: &entity.ShortURL{
					ID:  "10abcdef",
					URL: "http://full.url.com/test",
				},
			},
			want: want{
				code:    http.StatusTemporaryRedirect,
				headers: map[string]string{"Location": "http://full.url.com/test"},
			},
		},
		{
			name:         "Not found",
			request:      httptest.NewRequest(http.MethodGet, "/10abcdef", nil),
			usecaseParam: "10abcdef",
			usecaseRes: &usecaseResult{
				shortURL: nil,
				err:      entity.ErrShortURLNotFound,
			},
			want: want{
				code:    http.StatusBadRequest,
				headers: nil,
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {
			h.usecases.EXPECT().GetShortURL(test.usecaseParam).Return(test.usecaseRes.shortURL, test.usecaseRes.err).Once()

			// создаём новый Recorder
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(test.request, w)
			c.SetPath(":id")
			c.SetParamNames("id")
			c.SetParamValues(test.usecaseParam)

			if h.NoError(h.handlers.GetByShortURL(c)) {
				// проверяем код ответа
				h.Equal(test.want.code, w.Code)

				if test.want.headers != nil { // получаем и проверяем заголовки запроса
					for k, v := range test.want.headers {
						h.Equal(w.Header().Get(k), v)
					}
				}
			}
		})
	}
}

func (h *handlersSuite) TestCreateShortURL() {
	type usecaseResult struct {
		shortURL *entity.ShortURL
		err      error
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name         string
		request      *http.Request
		contentType  string
		usecaseParam string
		usecaseRes   *usecaseResult
		want         want
	}{
		{
			name:         "Success",
			request:      httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://full.url.com/test")),
			contentType:  "text/plain; charset=utf-8",
			usecaseParam: "http://full.url.com/test",
			usecaseRes: &usecaseResult{
				shortURL: &entity.ShortURL{
					ID:  "10abcdef",
					URL: "http://full.url.com/test",
				},
			},
			want: want{
				code:        http.StatusCreated,
				response:    h.config.BaseShortURL + "/10abcdef",
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name:        "Bad url",
			request:     httptest.NewRequest(http.MethodPost, "/", strings.NewReader("bad-url")),
			contentType: "text/plain; charset=utf-8",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name:        "Bad content type",
			request:     httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"url\": \"http://full.url.com/test\"}")),
			contentType: "application/json",
			want: want{
				code:     http.StatusUnsupportedMediaType,
				response: "",
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {
			if test.usecaseParam != "" && test.usecaseRes != nil {
				h.usecases.EXPECT().CreateShortURL(test.usecaseParam).Return(test.usecaseRes.shortURL, test.usecaseRes.err).Once()
			}

			// создаём новый Recorder
			w := httptest.NewRecorder()

			test.request.Header.Set("Content-Type", test.contentType)
			e := echo.New()
			c := e.NewContext(test.request, w)

			if h.NoError(h.handlers.CreateShortURL(c)) {
				// проверяем код ответа
				h.Equal(test.want.code, w.Code)

				if test.want.response != "" { // получаем и проверяем заголовки запроса
					resBody, err := io.ReadAll(w.Body)

					h.Require().NoError(err)
					h.Equal(test.want.response, string(resBody))
					h.Equal(test.want.contentType, w.Header().Get("Content-Type"))
				}
			}
		})
	}
}
