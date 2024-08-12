package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
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
