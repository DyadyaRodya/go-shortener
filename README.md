# go-shortener
Сервис сокращения URL


## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

### Локальный запуск автотестов

Установить `shortenertestbeta` [отсюда](https://github.com/Yandex-Practicum/go-autotests?tab=readme-ov-file#%D1%82%D1%80%D0%B5%D0%BA-%D1%81%D0%B5%D1%80%D0%B2%D0%B8%D1%81-%D1%81%D0%BE%D0%BA%D1%80%D0%B0%D1%89%D0%B5%D0%BD%D0%B8%D1%8F-url)
Установить `statictest` [отсюда](https://github.com/Yandex-Practicum/go-autotests?tab=readme-ov-file#%D1%82%D1%80%D0%B5%D0%BA-%D1%81%D0%B5%D1%80%D0%B2%D0%B8%D1%81-%D1%81%D0%BE%D0%BA%D1%80%D0%B0%D1%89%D0%B5%D0%BD%D0%B8%D1%8F-url)

Для статического теста
```shell
make lint
```

Для конкретной итерации выполняем 
```shell
make test-iter NUMBER=${ITER_NUMBER}
```

Для запуска всех тестов выполняем
```shell
make test-all
```

### Benchmark
Base:
```shell
cd internal/usecases/tests/
go test -bench=BenchmarkBatchCreateShortURLsWithDB -run=internal/usecases/tests/ . -benchmem -memprofile=../../../profiles/base.pprof
go tool pprof -http=":9090"  tests.test ../../../profiles/base.pprof
```
Result
```shell
cd internal/usecases/tests/
go test -bench=BenchmarkBatchCreateShortURLsWithDB -run=internal/usecases/tests/ . -benchmem -memprofile=../../../profiles/result.pprof
go tool pprof -http=":9090"  tests.test ../../../profiles/result.pprof

pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
```


### Дополнительные команды

Для генерации моков
```shell
make mock
```

Для запуска unit тестов выполняем
```shell
make tests
```
