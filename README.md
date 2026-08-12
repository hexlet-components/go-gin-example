# Go Gin Example

[![CI](https://github.com/hexlet-components/go-gin-example/actions/workflows/ci.yml/badge.svg)](https://github.com/hexlet-components/go-gin-example/actions/workflows/ci.yml)

## Зачем это нужно

REST API на [Gin](https://gin-gonic.com/): маршруты, обработчики, миграции
[goose](https://github.com/pressly/goose), сгенерированный
[sqlc](https://sqlc.dev/) слой доступа к базе и интеграционные тесты на
`httptest`.

Служит образцом раскладки Go-приложения, которое больше одного файла: команды
в *cmd*, обработчики в *handlers*, работа с базой в *db*. На него ссылается шаг
проекта «Сокращатель ссылок» как на пример, с которым можно сверяться.

## Requirements

* Go 1.26

## Commands

```bash
make install     # download dependencies
make db-migrate  # apply migrations
make dev         # run the API server at http://localhost:8080
make lint        # gofmt check and go vet
make test        # run tests
make help        # show all targets
```

База данных SQLite, файл создаётся миграциями рядом с проектом, ставить ничего
не нужно. Код доступа к базе генерируется из *db/queries* командой
`make db-generate`: `sqlc` объявлен в `go.mod` директивой `tool` и зовётся через
`go tool`, отдельная установка ему не нужна.

---

[![Hexlet Ltd. logo](https://raw.githubusercontent.com/Hexlet/assets/master/images/hexlet_logo128.png)](https://hexlet.io?utm_source=github&utm_medium=link&utm_campaign=go-gin-example)

This repository is created and maintained by the team and the community of Hexlet, an educational project. [Read more about Hexlet](https://hexlet.io?utm_source=github&utm_medium=link&utm_campaign=go-gin-example).

See most active contributors on [hexlet-friends](https://friends.hexlet.io/).
