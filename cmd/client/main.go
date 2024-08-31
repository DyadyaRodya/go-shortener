package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type RequestMaker interface {
	Make(long string) *http.Request
}

type JSONRequestMaker struct{}

func (JSONRequestMaker) Make(long string) *http.Request {
	endpoint := "http://localhost:8080/api/shorten"
	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	type Request struct {
		URL string `json:"url"`
	}

	body, err := json.Marshal(Request{URL: long})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "application/json")
	return request
}

type TextRequestMaker struct{}

func (TextRequestMaker) Make(long string) *http.Request {
	endpoint := "http://localhost:8080/"
	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(long))
	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "text/plain")
	return request
}

func main() {
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)

	// приглашение в консоли
	fmt.Println("JSON или Text?")
	// читаем строку из консоли
	reqType, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	reqType = strings.TrimSuffix(reqType, "\n")

	var requestMaker RequestMaker
	switch strings.ToLower(reqType) {
	case "text":
		requestMaker = TextRequestMaker{}
	case "json":
		requestMaker = JSONRequestMaker{}
	default:
		panic("unknown request type")
	}

	// приглашение в консоли
	fmt.Println("Введите длинный URL")
	// читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")
	// добавляем HTTP-клиент
	client := &http.Client{}

	// отправляем запрос и получаем ответ
	response, err := client.Do(requestMaker.Make(long))
	if err != nil {
		panic(err)
	}
	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// и печатаем его
	fmt.Println(string(body))
}
