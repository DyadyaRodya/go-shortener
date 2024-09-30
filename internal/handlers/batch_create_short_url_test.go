package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (h *handlersSuite) TestBatchCreateShortURLJSON() {
	type usecaseResult struct {
		responses []*usecasesdto.BatchCreateResponse
		err       error
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
		usecaseParam []*usecasesdto.BatchCreateRequest
		usecaseRes   *usecaseResult
		want         want
	}{
		{
			name: "Success",
			request: httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(
				`[
					  {
						"correlation_id":"1",
						"original_url":"http://full.url.com/test"
					  }
					]`)),
			contentType: "application/json; charset=utf-8",
			usecaseParam: []*usecasesdto.BatchCreateRequest{
				&usecasesdto.BatchCreateRequest{
					CorrelationID: "1",
					OriginalURL:   "http://full.url.com/test",
				},
			},
			usecaseRes: &usecaseResult{
				responses: []*usecasesdto.BatchCreateResponse{
					&usecasesdto.BatchCreateResponse{
						CorrelationID: "1",
						ShortURL: &entity.ShortURL{
							ID:  "10abcdef",
							URL: "http://full.url.com/test",
						},
					},
				},
			},
			want: want{
				code: http.StatusCreated,
				response: `[
							 {
							   "correlation_id":"1",
							   "short_url":"` + h.config.BaseShortURL + "/10abcdef" + `"
							 }
						   ]`,
				contentType: "application/json",
			},
		},
		{
			name:        "Bad_url",
			request:     httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(`{"url":"bad-url"}`)),
			contentType: "application/json; charset=utf-8",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name:        "Not_json",
			request:     httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader("http://full.url.com/test")),
			contentType: "text/plain",
			want: want{
				code:     http.StatusBadRequest,
				response: "",
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
			if test.usecaseParam != nil && test.usecaseRes != nil {
				h.usecases.EXPECT().BatchCreateShortURLs(ctx, test.usecaseParam, userUUID).Return(test.usecaseRes.responses, test.usecaseRes.err).Once()
			}

			if h.NoError(h.handlers.BatchCreateShortURLJSON(c)) {
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
