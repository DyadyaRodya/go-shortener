package handlers

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (h *handlersSuite) TestGetShortURL() {
	type usecaseResult struct {
		shortURL *entity.ShortURL
		err      error
	}
	type want struct {
		headers map[string]string
		code    int
	}
	tests := []struct {
		name         string
		usecaseParam string
		request      *http.Request
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
		{
			name:         "Deleted",
			request:      httptest.NewRequest(http.MethodGet, "/10abcdef", nil),
			usecaseParam: "10abcdef",
			usecaseRes: &usecaseResult{
				shortURL: nil,
				err:      entity.ErrShortURLDeleted,
			},
			want: want{
				code:    http.StatusGone,
				headers: nil,
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(test.request, w)
			c.SetPath(":id")
			c.SetParamNames("id")
			c.SetParamValues(test.usecaseParam)

			ctx := c.Request().Context()
			h.usecases.EXPECT().GetShortURL(ctx, test.usecaseParam).Return(test.usecaseRes.shortURL, test.usecaseRes.err).Once()

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
