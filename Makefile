LOCAL_BIN:=$(CURDIR)/bin

.PHONY: install_bin
install_bin: # install binary dependencies
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go mod tidy
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@latest
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@latest
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/godoc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: install
install: install_bin

.PHONY:
mockery:
	$(LOCAL_BIN)/mockery --name $(name) --dir $(dir) --output $(dir)/mocks

.PHONY:
mock:
	make mockery name=Usecases dir=./internal/handlers

	make mockery name=Transaction dir=./internal/usecases
	make mockery name=URLStorage dir=./internal/usecases
	make mockery name=IDGenerator dir=./internal/usecases

.PHONY: lint
lint: # run statictest
	$(LOCAL_BIN)/goimports -local "github.com/DyadyaRodya/go-shortener" -w cmd internal pkg
	go vet -vettool=/usr/bin/statictest ./...
	go build -o cmd/staticlint/main cmd/staticlint/main.go && go vet -vettool=cmd/staticlint/main ./...

.PHONY: tests
tests: # run unit tests
	go test -race -coverprofile=coverage.out ./...

.PHONY: test-iter
test-iter: # run test for iteration
	if [ ${NUMBER} -ge 1 ]; then shortenertestbeta -test.v -test.run=^TestIteration1$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 2 ]; then shortenertestbeta -test.v -test.run=^TestIteration2$$ -source-path=.; fi; \
	if [ ${NUMBER} -ge 3 ]; then shortenertestbeta -test.v -test.run=^TestIteration3$$ -source-path=.; fi; \
	if [ ${NUMBER} -ge 4 ]; then shortenertestbeta -test.v -test.run=^TestIteration4$$ \
									 -binary-path=cmd/shortener/shortener \
									 -server-port=8888; fi; \
	if [ ${NUMBER} -ge 5 ]; then shortenertestbeta -test.v -test.run=^TestIteration5$$ \
								  	 -binary-path=cmd/shortener/shortener \
									 -server-port=8888; fi; \
	if [ ${NUMBER} -ge 6 ]; then shortenertestbeta -test.v -test.run=^TestIteration6$$ \
              						-source-path=.; fi; \
	if [ ${NUMBER} -ge 7 ]; then shortenertestbeta -test.v -test.run=^TestIteration7$$ \
									-binary-path=cmd/shortener/shortener \
									-source-path=.; fi; \
	if [ ${NUMBER} -ge 8 ]; then shortenertestbeta -test.v -test.run=^TestIteration8$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 9 ]; then shortenertestbeta -test.v -test.run=^TestIteration9$$ \
								 	 -binary-path=cmd/shortener/shortener \
								 	 -source-path=. \
								  	 -file-storage-path=/tmp/go-shortener; fi; \
	if [ ${NUMBER} -ge 10 ]; then shortenertestbeta -test.v -test.run=^TestIteration10$$  \
									 -binary-path=cmd/shortener/shortener \
									 -source-path=. \
									 -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'; fi; \
	if [ ${NUMBER} -ge 11 ]; then shortenertestbeta -test.v -test.run=^TestIteration11$$ \
									 -binary-path=cmd/shortener/shortener \
									 -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'; fi; \
	if [ ${NUMBER} -ge 12 ]; then shortenertestbeta -test.v -test.run=^TestIteration12$$ \
                                     -binary-path=cmd/shortener/shortener \
                                     -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'; fi; \
	if [ ${NUMBER} -ge 13 ]; then shortenertestbeta -test.v -test.run=^TestIteration13$$ \
                                     -binary-path=cmd/shortener/shortener \
                                     -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'; fi; \
	if [ ${NUMBER} -ge 14 ]; then shortenertestbeta -test.v -test.run=^TestIteration14$$ \
                                     -binary-path=cmd/shortener/shortener \
                                     -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'; fi; \
	if [ ${NUMBER} -ge 15 ]; then shortenertestbeta -test.v -test.run=^TestIteration15$$ \
                                     -binary-path=cmd/shortener/shortener \
                                     -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'; fi; \
	if [ ${NUMBER} -ge 16 ]; then shortenertestbeta -test.v -test.run=^TestIteration16$$ -source-path=. ; fi; \
	if [ ${NUMBER} -ge 17 ]; then shortenertestbeta -test.v -test.run=^TestIteration17$$ -source-path=. ; fi; \
	if [ ${NUMBER} -ge 18 ]; then shortenertestbeta -test.v -test.run=^TestIteration18$$ -source-path=. ; fi

.PHONY: test-all
test-all: # run test for all iterations
	shortenertestbeta -test.v -test.run=^TestIteration -binary-path=cmd/shortener/shortener


.PHONY: swagger
swagger: # generate swagger
	$(LOCAL_BIN)/swag init --pd --dir internal/handlers --generalInfo handlers.go
