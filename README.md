# Webserver

Простой веб-сервер на Go для демонстрации работы с JSON и HTML-шаблонами.

## Структура проекта

```
src/
├── main.go          # Точка входа
├── http/            # HTTP-обработчики
│   └── server.go
├── static/          # Статические файлы и шаблоны
└── test/            # Тесты
    └── server_test.go
```

## API endpoints

| Endpoint | Описание |
|----------|----------|
| `GET /` | Главная страница (статика) |
| `GET /api/items` | JSON API с данными |
| `GET /items` | HTML-страница со списком |
| `GET /itemsjson` | HTML-страница с JSON внутри |

## Запуск

```bash
cd src
go run main.go
```

Сервер запустится на `http://localhost:8080`

## Тестирование

```bash
cd src
go test ./test/ -v
```

### Покрытие тестов
- ✅ JSON API endpoint
- ✅ HTML-шаблоны
- ✅ Обработка ошибок
- ✅ Парсинг JSON из HTML

## Особенности

- **Fallback для шаблонов**: если файлы шаблонов недоступны (например, в тестах), используются встроенные константы
- **Без внешних зависимостей** для работы (только `testify` для тестов)
- **Docker-ready**: есть Dockerfile для контейнеризации

## Технологии

- Go 1.25.3
- Стандартная библиотека `net/http`
- Шаблонизатор `text/template`
- Тесты: `github.com/stretchr/testify`
