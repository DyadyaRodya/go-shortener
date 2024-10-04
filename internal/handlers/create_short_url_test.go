package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

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
			name:        "Bad_url",
			request:     httptest.NewRequest(http.MethodPost, "/", strings.NewReader("bad-url")),
			contentType: "text/plain; charset=utf-8",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name:        "Bad_content_type",
			request:     httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"url\": \"http://full.url.com/test\"}")),
			contentType: "application/json",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name:         "ShortURL_exists",
			request:      httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://full.url.com/test")),
			contentType:  "text/plain; charset=utf-8",
			usecaseParam: "http://full.url.com/test",
			usecaseRes: &usecaseResult{
				shortURL: &entity.ShortURL{
					ID:  "10abcdef",
					URL: "http://full.url.com/test",
				},
				err: entity.ErrShortURLExists,
			},
			want: want{
				code:        http.StatusConflict,
				response:    h.config.BaseShortURL + "/10abcdef",
				contentType: "text/plain; charset=UTF-8",
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			test.request.Header.Set("Content-Type", test.contentType)
			e := echo.New()
			c := e.NewContext(test.request, w)

			userUUID := gofakeit.UUID()
			c.Set("userUUID", userUUID)

			ctx := c.Request().Context()
			if test.usecaseParam != "" && test.usecaseRes != nil {
				h.usecases.EXPECT().CreateShortURL(ctx, test.usecaseParam, userUUID).Return(test.usecaseRes.shortURL, test.usecaseRes.err).Once()
			}

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

func (h *handlersSuite) TestCreateShortURLJSON() {
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
			request:      httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"http://full.url.com/test"}`)),
			contentType:  "application/json; charset=utf-8",
			usecaseParam: "http://full.url.com/test",
			usecaseRes: &usecaseResult{
				shortURL: &entity.ShortURL{
					ID:  "10abcdef",
					URL: "http://full.url.com/test",
				},
			},
			want: want{
				code:        http.StatusCreated,
				response:    `{"result":"` + h.config.BaseShortURL + "/10abcdef" + `"}`,
				contentType: "application/json",
			},
		},
		{
			name:        "Bad_url",
			request:     httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"bad-url"}`)),
			contentType: "application/json; charset=utf-8",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name:        "Not_json",
			request:     httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader("http://full.url.com/test")),
			contentType: "text/plain",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name:         "ShortURL_exists",
			request:      httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"http://full.url.com/test"}`)),
			contentType:  "application/json; charset=utf-8",
			usecaseParam: "http://full.url.com/test",
			usecaseRes: &usecaseResult{
				shortURL: &entity.ShortURL{
					ID:  "10abcdef",
					URL: "http://full.url.com/test",
				},
				err: entity.ErrShortURLExists,
			},
			want: want{
				code:        http.StatusConflict,
				response:    `{"result":"` + h.config.BaseShortURL + "/10abcdef" + `"}`,
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			test.request.Header.Set("Content-Type", test.contentType)
			e := echo.New()
			c := e.NewContext(test.request, w)

			userUUID := gofakeit.UUID()
			c.Set("userUUID", userUUID)

			ctx := c.Request().Context()
			if test.usecaseParam != "" && test.usecaseRes != nil {
				h.usecases.EXPECT().CreateShortURL(ctx, test.usecaseParam, userUUID).Return(test.usecaseRes.shortURL, test.usecaseRes.err).Once()
			}

			if h.NoError(h.handlers.CreateShortURLJSON(c)) {
				// проверяем код ответа
				h.Equal(test.want.code, w.Code)

				if test.want.response != "" { // получаем и проверяем заголовки запроса
					resBody, err := io.ReadAll(w.Body)

					h.Require().NoError(err)
					h.JSONEq(test.want.response, string(resBody))
					h.Equal(test.want.contentType, w.Header().Get("Content-Type"))
				}
			}
		})
	}
}
