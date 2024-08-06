package handlers

import (
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	handlersMocks "github.com/DyadyaRodya/go-shortener/internal/handlers/mocks"
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

			h.handlers.RootHandler(w, test.request) // want to protect API so root handler for now

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			h.Equal(test.want.code, res.StatusCode)

			if test.want.headers != nil { // получаем и проверяем заголовки запроса
				for k, v := range test.want.headers {
					h.Equal(res.Header.Get(k), v)
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
				contentType: "text/plain",
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
			h.handlers.RootHandler(w, test.request) // want to protect API so root handler for now

			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			h.Equal(test.want.code, res.StatusCode)

			if test.want.response != "" { // получаем и проверяем заголовки запроса
				resBody, err := io.ReadAll(res.Body)

				h.Require().NoError(err)
				h.Equal(test.want.response, string(resBody))
				h.Equal(test.want.contentType, res.Header.Get("Content-Type"))
			}

		})
	}
}

func (h *handlersSuite) TestRootHandler_MethodNotAllowed() {
	methods := []string{
		http.MethodDelete,
		http.MethodPut,
		http.MethodHead,
		http.MethodPatch,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range methods {
		h.Run(method, func() {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			// создаем новый Request
			r := httptest.NewRequest(method, "/", nil)

			h.handlers.RootHandler(w, r)
			res := w.Result()
			defer res.Body.Close()

			h.Equal(http.StatusMethodNotAllowed, res.StatusCode)
		})
	}
}
