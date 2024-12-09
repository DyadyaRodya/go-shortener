package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (h *handlersSuite) TestGetUserShortURLs() {
	type usecaseResult struct {
		err       error
		shortURLs []*entity.ShortURL
	}
	type want struct {
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		name       string
		userUUID   string
		request    *http.Request
		usecaseRes *usecaseResult
		want       want
		authorized bool
	}{
		{
			name:       "Success",
			request:    httptest.NewRequest(http.MethodGet, "/api/user/urls", nil),
			authorized: true,
			userUUID:   gofakeit.UUID(),
			usecaseRes: &usecaseResult{
				shortURLs: []*entity.ShortURL{
					&entity.ShortURL{
						ID:  "10abcdef",
						URL: "http://full.url.com/test",
					},
				},
			},
			want: want{
				code:        http.StatusOK,
				response:    `[{"short_url": "` + h.config.BaseShortURL + `/10abcdef", "original_url": "http://full.url.com/test"}]`,
				contentType: "application/json",
			},
		},
		{
			name:       "Not_found",
			request:    httptest.NewRequest(http.MethodGet, "/api/user/urls", nil),
			authorized: true,
			userUUID:   gofakeit.UUID(),
			usecaseRes: &usecaseResult{
				shortURLs: []*entity.ShortURL{},
				err:       nil,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name:       "Not_authorized",
			request:    httptest.NewRequest(http.MethodGet, "/api/user/urls", nil),
			authorized: false,
			userUUID:   "",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(test.request, w)

			c.Set("authorized", test.authorized)
			c.Set("userUUID", test.userUUID)

			ctx := c.Request().Context()
			if test.authorized && test.userUUID != "" {
				h.usecases.EXPECT().GetUserShortURLs(ctx, test.userUUID).Return(test.usecaseRes.shortURLs, test.usecaseRes.err).Once()
			}

			if h.NoError(h.handlers.GetUserShortURLs(c)) {
				// проверяем код ответа
				h.Equal(test.want.code, w.Code)

				if test.want.response != "" { // получаем и проверяем ответ и заголовки запроса
					resBody, err := io.ReadAll(w.Body)

					h.Require().NoError(err)
					h.JSONEq(test.want.response, string(resBody))
					h.Equal(test.want.contentType, w.Header().Get("Content-Type"))
				}
			}
		})
	}
}
