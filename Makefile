LOCAL_BIN:=$(CURDIR)/bin

.PHONY: install_bin
install_bin: # install binary dependencies
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go mod tidy
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@latest

.PHONY: install
install: install_bin

.PHONY:
mockery:
	$(LOCAL_BIN)/mockery --name $(name) --dir $(dir) --output $(dir)/mocks

.PHONY:
mock:
	make mockery name=Usecases dir=./internal/handlers

	make mockery name=URLStorage dir=./internal/usecases
	make mockery name=IDGenerator dir=./internal/usecases

.PHONY: tests
tests: # run unit tests
	go test -race -coverprofile=coverage.out ./...

.PHONY: test-iter
test-iter: # run test for iteration
	if [ ${NUMBER} -ge 1 ]; then shortenertest -test.v -test.run=^TestIteration1$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 2 ]; then shortenertest -test.v -test.run=^TestIteration2$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 3 ]; then shortenertest -test.v -test.run=^TestIteratio3$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 4 ]; then shortenertest -test.v -test.run=^TestIteration4$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 5 ]; then shortenertest -test.v -test.run=^TestIteration5$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 6 ]; then shortenertest -test.v -test.run=^TestIteration6$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 7 ]; then shortenertest -test.v -test.run=^TestIteration7$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 8 ]; then shortenertest -test.v -test.run=^TestIteration8$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 9 ]; then shortenertest -test.v -test.run=^TestIteration9$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 10 ]; then shortenertest -test.v -test.run=^TestIteration10$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 11 ]; then shortenertest -test.v -test.run=^TestIteration11$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 12 ]; then shortenertest -test.v -test.run=^TestIteration12$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 13 ]; then shortenertest -test.v -test.run=^TestIteration13$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 14 ]; then shortenertest -test.v -test.run=^TestIteration14$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 15 ]; then shortenertest -test.v -test.run=^TestIteration15$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 16 ]; then shortenertest -test.v -test.run=^TestIteration16$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 17 ]; then shortenertest -test.v -test.run=^TestIteration17$$ -binary-path=cmd/shortener/shortener; fi; \
	if [ ${NUMBER} -ge 18 ]; then shortenertest -test.v -test.run=^TestIteration18$$ -binary-path=cmd/shortener/shortener; fi

.PHONY: test-all
test-all: # run test for all iterations
	shortenertest -test.v -test.run=^TestIteration -binary-path=cmd/shortener/shortener