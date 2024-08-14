package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
)

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
			name:         "Not_found",
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
