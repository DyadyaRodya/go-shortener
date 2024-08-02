.PHONY: test
test: # run test for iteration
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
test-all: # run test for iteration
	shortenertest -test.v -test.run=^TestIteration -binary-path=cmd/shortener/shortener