# Docker и CI для API-автотеста

Проект запускает GraphQL query, mutation и негативный тест к `https://graphql.postman-echo.com/graphql`

## Локальный запуск

Сборка Docker image:

```bash
docker build -t api-test .
```

Запуск теста в контейнере:

```bash
docker run --rm api-test
```

Контейнер автоматически выполняет `go test -count=1 -v ./...`

## CI

Конфигурация находится в [`.github/workflows/ci.yml`](.github/workflows/ci.yml)

GitHub Actions запускается при каждом `push` и вручную через `workflow_dispatch`

[Успешный pipeline](https://github.com/ScreoX/api-testing-ci/actions/workflows/ci.yml)
