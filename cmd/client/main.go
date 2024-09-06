package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
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

type JSONRequestMaker struct {
	gzip bool
}

func (m JSONRequestMaker) Make(long string) *http.Request {
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

	var buf bytes.Buffer
	if m.gzip {
		g := gzip.NewWriter(&buf)
		if _, err = g.Write(body); err != nil {
			panic(err)

		}
		if err = g.Close(); err != nil {
			panic(err)
		}
	} else {
		if _, err = buf.Write(body); err != nil {
			panic(err)
		}
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, &buf)

	if err != nil {
		panic(err)
	}
	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "application/json")
	if m.gzip {
		request.Header.Add("Content-Encoding", "gzip")
	} else {
		request.Header.Del("Content-Encoding")
	}
	return request
}

type TextRequestMaker struct {
	gzip bool
}

func (m TextRequestMaker) Make(long string) *http.Request {
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

	// приглашение в консоли
	fmt.Println("gzip?")
	// читаем строку из консоли
	compress, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	useGZIP := strings.ToLower(strings.TrimSuffix(compress, "\n")) == "yes"

	var requestMaker RequestMaker
	switch strings.ToLower(reqType) {
	case "text":
		requestMaker = TextRequestMaker{useGZIP}
	case "json":
		requestMaker = JSONRequestMaker{useGZIP}
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

	//gzip
	tr := &http.Transport{
		DisableCompression: !useGZIP,
	}
	// добавляем HTTP-клиент
	client := &http.Client{
		Transport: tr,
	}

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
