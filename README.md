# loglint

`loglint` — кастомный линтер для Go, который проверяет корректность лог-сообщений в вызовах `log/slog` и `go.uber.org/zap`.

Линтер построен на базе `golang.org/x/tools/go/analysis` и может запускаться:

- как standalone-анализатор (`go run ./cmd/loglint ./...`);
- как кастомный линтер через `golangci-lint` (module plugin mode).

## Что проверяет линтер

`loglint` валидирует 4 правила:

1. Сообщение должно начинаться со строчной буквы.
2. Сообщение должно быть только на английском языке.
3. Сообщение не должно содержать спецсимволы, шумную пунктуацию или эмодзи.
4. Сообщение не должно содержать потенциально чувствительные данные.

Примеры некорректных сообщений:

```go
slog.Info("Starting server")
slog.Info("запуск сервера")
logger.Error("connection failed!!!")
logger.Info("token: abc123")
```

## Поддерживаемые логгеры

- `log/slog`
- `go.uber.org/zap`

## Структура проекта

```text
cmd/loglint              точка входа standalone-режима (singlechecker)
pkg/analyzer             основной анализатор
pkg/logger               матчинг вызовов slog/zap
pkg/rules                правила валидации сообщений
plugin                   регистрация плагина для golangci-lint
pkg/analyzer/testdata    фикстуры analysistest для analyzer
pkg/logger/testdata      фикстуры analysistest для matcher
```

## Требования

- Go `1.22+` (в `go.mod` сейчас `go 1.24.4`);
- `golangci-lint` (для custom-сборки и запуска плагина).

## Установка

Установить `golangci-lint`:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Если `golangci-lint` не найден, добавьте Go bin в `PATH`.

PowerShell (текущая сессия):

```powershell
$env:Path += ";$env:USERPROFILE\\go\\bin"
```

Git Bash:

```bash
export PATH="$PATH:/c/Users/<your-user>/go/bin"
```

## Локальный запуск (standalone)

```bash
go run ./cmd/loglint ./...
```

Проверить только демонстрационный файл с логами:

```powershell
$env:GOFLAGS='-tags=manualcheck'
go run ./cmd/loglint ./examples/manualcheck
```

```bash
GOFLAGS='-tags=manualcheck' go run ./cmd/loglint ./examples/manualcheck
```

## Запуск тестов

```bash
go test ./...
```

## Запуск через golangci-lint (module plugin)

В репозитории уже есть конфиги:

- `.custom-gcl.yml` — сборка custom-бинаря `golangci-lint` с вашим плагином;
- `.golangci.yml` — включение и запуск линтера `loglint`.

Собрать custom-бинарь:

```bash
golangci-lint custom -v
```

Запустить проверку:

PowerShell:

```powershell
.\custom-gcl.exe run ./...
```

Git Bash:

```bash
./custom-gcl.exe run ./...
```

## CI

Workflow [`.github/workflows/ci.yml`](.github/workflows/ci.yml) запускает:

- `go test ./...`
- `go run ./cmd/loglint ./...`
- manual build of custom `golangci-lint` binary (module plugin mode)
- `./custom-gcl run ./...`

## Частые проблемы

- `golangci-lint: command not found`
  - `golangci-lint` не установлен или путь к Go bin не добавлен в `PATH`.
- `plugin "loglint" not found`
  - пересоберите custom-бинарь из корня проекта: `golangci-lint custom -v`.
- `bash: .custom-gcl.exe: command not found`
  - в Git Bash используйте `./custom-gcl.exe`, а не `.custom-gcl.exe`.
