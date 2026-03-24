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
(cd cmd/gitfame && go build .)
```

## Установка

```bash
go install ./cmd/gitfame/...
```

## Запуск

```bash
gitfame --repository=. --revision=HEAD --order-by=lines --format=tabular
```

Пример tabular-вывода:

```text
Name  Lines Commits Files
Alice 10    2       2
Bob   3     1       1
```

Пример JSON-вывода:

```json
[{"name":"Alice","lines":10,"commits":2,"files":2},{"name":"Bob","lines":3,"commits":1,"files":1}]
```

Пример с фильтрами:

```bash
gitfame --repository=. --extensions='.go,.md' --languages='go,markdown' --exclude='vendor/*'
```

## Как считается статистика

- Список файлов берётся из Git tree выбранной ревизии (`git ls-tree`), поэтому учитываются только tracked-файлы в состоянии этой ревизии.
- Attribution строк считается через `git blame --porcelain`.
- Для пустых файлов берётся последний коммит, изменявший файл (`git log -n 1`): увеличиваются `Commits` и `Files`, а `Lines` остаётся `0`.

## Тесты

Все тесты:

```bash
go test ./...
```

Интеграционные тесты:

```bash
go test -v ./test/integration/...
```

## Важные детали

- Неизвестные языки в `--languages` не ограничивают выборку и выводят warning в `stderr`.
- Невалидные значения флагов (включая невалидные glob-паттерны) завершают программу с ненулевым кодом.
- Утилита не изменяет состояние репозитория.
