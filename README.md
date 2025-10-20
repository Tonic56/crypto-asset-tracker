# crypto-asset-tracker

Минимальный проект для сбора и агрегации торговых данных (WIP).

Как собрать:

```bash
go build ./...
```

Как запустить (локально):

1. Откройте `cmd/main.go` и настройте URL для websocket (если нужно).
2. Запустите:

```bash
go run ./cmd
```

Что есть:
- `internal/websocket` — WebSocket клиент
- `internal/processor` — парсер и конвертеры сообщений
- `internal/aggregator` — заготовка для агрегации свечей/окон
- `internal/kafka` — (заглушки для продюсера/батчера/партиционирования)

Дальше:
- Реализовать логику Aggregator.Start и processIncoming
- Добавить обработку graceful shutdown
- Добавить тесты и CI
