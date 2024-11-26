package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"

	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

func (h *handlersSuite) TestDeleteUserShortURLs() {

	type want struct {
		code int
	}
	tests := []struct {
		name        string
		request     *http.Request
		contentType string
		authorized  bool
		userUUID    string
		want        want
	}{
		{
			name: "Success",
			request: httptest.NewRequest(http.MethodDelete, "/api/shorten/batch",
				strings.NewReader(`["10abcdef"]`)),
			contentType: "application/json; charset=utf-8",
			authorized:  true,
			userUUID:    gofakeit.UUID(),
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: "Bad_id",
			request: httptest.NewRequest(http.MethodDelete, "/api/shorten/batch",
				strings.NewReader(`["blablablabla"]`)),
			contentType: "application/json; charset=utf-8",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Not_json",
			request: httptest.NewRequest(http.MethodDelete, "/api/shorten/batch",
				strings.NewReader("10abcdef")),
			contentType: "text/plain",
			want: want{
				code: http.StatusBadRequest,
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

			c.Set("authorized", test.authorized)
			c.Set("userUUID", test.userUUID)

			if h.NoError(h.handlers.DeleteUserShortURLs(c)) {
				// test code
				h.Equal(test.want.code, w.Code)

				if test.want.code < 400 {
					// test DelChan has request
					select {
					case data := <-h.handlers.DelChan:
						h.Equal(usecasesdto.DeleteUserShortURLsRequest{
							UserUUID: test.userUUID, ShortURLUUIDs: []string{"10abcdef"},
						}, *data)
					case <-time.After(1 * time.Second): // to let data appear in chan
						h.Fail("No data in h.handlers.DelChan")
					}
				} else {
					select {
					case data := <-h.handlers.DelChan:
						h.Failf("data in h.handlers.DelChan, but not expected",
							"data in h.handlers.DelChan, but not expected: %+v", *data)
					case <-time.After(1 * time.Second): // to let data appear in chan
						// ok
					}
				}
			}
		})
	}
}
