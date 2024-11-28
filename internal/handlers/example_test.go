package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var cookie = http.Cookie{
	Name:     "Auth",
	Value:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	HttpOnly: true,
	Path:     "/",
}

func ExampleHandlers_BatchCreateShortURLJSON() {
	httpClient := &http.Client{}
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/shorten/batch",
		strings.NewReader(`[{"correlation_id":"1","original_url":"http://full.url.com/test"}]`),
	)

	req.AddCookie(&cookie)

	resp, _ := httpClient.Do(req)

	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func ExampleHandlers_CreateShortURLJSON() {
	httpClient := &http.Client{}
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/shorten",
		strings.NewReader(`{"url":"http://full.url.com/test"}`),
	)

	req.AddCookie(&cookie)

	resp, _ := httpClient.Do(req)

	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func ExampleHandlers_DeleteUserShortURLs() {
	httpClient := &http.Client{}
	req, _ := http.NewRequest(
		http.MethodDelete,
		"http://localhost:8080/api/user/urls",
		strings.NewReader(`["10abcdef"]`),
	)

	req.AddCookie(&cookie)

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
}

func ExampleHandlers_GetByShortURL() {
	httpClient := &http.Client{}
	req, _ := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/10abcdef",
		strings.NewReader(""),
	)

	req.AddCookie(&cookie)

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Location"))
}

func ExampleHandlers_GetUserShortURLs() {
	httpClient := &http.Client{}
	req, _ := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/api/user/urls",
		strings.NewReader(""),
	)

	req.AddCookie(&cookie)

	resp, _ := httpClient.Do(req)

	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
