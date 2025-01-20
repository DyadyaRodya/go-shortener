package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

func (h *handlersSuite) TestGetStats() {
	type usecaseResult struct {
		stats *usecasesdto.StatsResponse
		err   error
	}
	type want struct {
		response    string
		contentType string
		code        int
	}
	tests := []struct {
		name       string
		request    *http.Request
		usecaseRes *usecaseResult
		want       want
	}{
		{
			name:    "Success",
			request: httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil),
			usecaseRes: &usecaseResult{
				stats: &usecasesdto.StatsResponse{
					URLs:  100500,
					Users: 9999,
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				response:    `{"urls": 100500, "users": 9999}`,
			},
		},
		{
			name:    "InternalError",
			request: httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil),
			usecaseRes: &usecaseResult{
				err: errors.New("some error"),
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		h.Run(test.name, func() {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(test.request, w)

			ctx := c.Request().Context()
			h.usecases.EXPECT().GetStats(ctx).Return(test.usecaseRes.stats, test.usecaseRes.err).Once()

			if h.NoError(h.handlers.GetStats(c)) {
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
