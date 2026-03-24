# gitfame

`gitfame` — CLI-утилита для подсчёта статистики по авторам Git-репозитория на выбранной ревизии.

## Что считает

- количество строк (`Lines`)
- количество уникальных коммитов (`Commits`)
- количество уникальных файлов (`Files`)

## Возможности

- расчёт по `author` (по умолчанию) или `committer` (`--use-committer`)
- фильтрация по расширениям (`--extensions`)
- фильтрация по языкам (`--languages`)
- исключение файлов через glob (`--exclude`)
- ограничение выборки через glob (`--restrict-to`)
- форматы вывода: `tabular`, `csv`, `json`, `json-lines`

## Сборка

```bash
(cd gitfame/cmd/gitfame && go build .)
```

## Установка

```bash
go install ./gitfame/cmd/gitfame/...
```

## Запуск

```bash
gitfame --repository=. --revision=HEAD --order-by=lines --format=tabular
```

Пример с фильтрами:

```bash
gitfame --repository=. --extensions='.go,.md' --languages='go,markdown' --exclude='vendor/*'
```

## Тесты

Все тесты:

```bash
go test ./...
```

Интеграционные тесты:

```bash
go test -v ./gitfame/test/integration/...
```

## Важные детали

- Неизвестные языки в `--languages` не ограничивают выборку и выводят warning в `stderr`.
- Невалидные значения флагов (включая невалидные glob-паттерны) завершают программу с ненулевым кодом.
- Утилита не изменяет состояние репозитория.
