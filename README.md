# Go-Formist

Go-Formist - библиотека для создания админ-панелей с формами на Go. Библиотека предоставляет простой и гибкий API для создания форм, таблиц и кастомных страниц с автоматической генерацией JSON Schema и UI Schema.

## Особенности

- 🚀 **Простой API** - Fluent interface для быстрого создания форм
- 📋 **Множество типов полей** - Text, Email, Password, Number, Select, Checkbox, Textarea, Date, File, Table
- 🔧 **Автогенерация из структур** - Создание форм из Go структур с помощью тегов
- 📊 **Встроенные таблицы** - Поддержка таблиц с сортировкой, фильтрацией и пагинацией
- ✅ **Валидация** - Встроенная валидация полей с кастомными правилами
- 🌐 **JSON Schema** - Автоматическая генерация JSON Schema и UI Schema
- 🔒 **Авторизация** - Встроенная поддержка авторизации
- 🌍 **CORS** - Настраиваемая поддержка CORS
- 🎨 **Кастомные страницы** - Возможность добавления собственных HTML страниц

## Установка

```bash
go get github.com/koteyye/go-formist
```

## Быстрый старт

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/koteyye/go-formist"
    "github.com/koteyye/go-formist/types"
)

func main() {
    // Создаем админ-панель
    admin := formist.New()
    admin.SetTitle("Моя Админ-панель")

    // Создаем простую форму
    userForm := formist.NewForm("user", "Пользователь").
        AddTextField("name", "Имя").
        AddEmailField("email", "Email").
        AddPasswordField("password", "Пароль").
        OnPost(func(data map[string]interface{}) (interface{}, error) {
            fmt.Printf("Данные: %+v\n", data)
            return map[string]string{"message": "Пользователь создан"}, nil
        }).
        Build()

    // Регистрируем форму
    admin.RegisterForm(userForm)

    // Запускаем сервер
    log.Fatal(http.ListenAndServe(":8080", admin.Handler()))
}
```

## Типы полей

### Базовые поля

```go
form := formist.NewForm("example", "Пример").
    AddTextField("name", "Имя").
    AddEmailField("email", "Email").
    AddPasswordField("password", "Пароль").
    AddNumberField("age", "Возраст").
    AddTextareaField("bio", "Биография").
    AddDateField("birthday", "День рождения").
    AddFileField("avatar", "Аватар").
    AddCheckboxField("active", "Активен").
    AddHiddenField("id", 123).
    Build()
```

### Select поля

```go
options := []types.SelectOption{
    formist.SelectOption("admin", "Администратор"),
    formist.SelectOption("user", "Пользователь"),
}

form := formist.NewForm("example", "Пример").
    AddSelectField("role", "Роль", options).
    AddMultiSelectField("permissions", "Права", options).
    Build()
```

### Таблицы

```go
formBuilder := formist.NewForm("orders", "Заказы")

tableField := formBuilder.AddTableField("orders_table", "Список заказов").
    AddTextColumn("id", "ID").WithSortable().
    AddTextColumn("customer", "Клиент").WithFilterable().
    AddEmailColumn("email", "Email").
    AddNumberColumn("amount", "Сумма").WithSortable().
    AddSelectColumn("status", "Статус", statusOptions).WithFilterable().
    WithPagination(true).
    WithPageSize(20).
    OnGet(func(page, limit int, filters map[string]interface{}) (types.TableData, error) {
        // Логика получения данных
        return types.TableData{
            Rows:  rows,
            Total: total,
            Page:  page,
            Limit: limit,
        }, nil
    })

form := tableField.Build(formBuilder).Build()
```

## Создание форм из структур

```go
type User struct {
    Name     string `form:"name" label:"Имя" required:"true"`
    Email    string `form:"email" label:"Email" type:"email"`
    Age      int    `form:"age" label:"Возраст" type:"number"`
    Active   bool   `form:"active" label:"Активен"`
    Bio      string `form:"bio" label:"Биография" type:"textarea"`
}

form := formist.FromStruct("user", "Пользователь", User{}).
    OnPost(func(data map[string]interface{}) (interface{}, error) {
        // Обработка данных
        return nil, nil
    }).
    Build()
```

### Поддерживаемые теги

- `form:"field_name"` - имя поля
- `label:"Field Label"` - метка поля
- `type:"field_type"` - тип поля (email, password, textarea, select, etc.)
- `required:"true"` - обязательное поле

## Валидация

```go
form := formist.NewForm("user", "Пользователь").
    AddField(types.Field{
        Name:     "email",
        Type:     types.FieldTypeEmail,
        Label:    "Email",
        Required: true,
        Validation: []types.ValidationRule{
            formist.ValidationRule("email", nil, "Некорректный email"),
            formist.ValidationRule("minLength", 5, "Минимум 5 символов"),
        },
    }).
    Build()
```

### Типы валидации

- `email` - валидация email адреса
- `min` / `max` - минимальное/максимальное значение для чисел
- `minLength` / `maxLength` - минимальная/максимальная длина строки
- `pattern` - валидация по регулярному выражению

## Кастомные страницы

```go
page := formist.NewPage("dashboard", "Панель управления").
    WithContent(`
        <h1>Добро пожаловать!</h1>
        <p>Это кастомная страница</p>
    `).
    Build()

admin.RegisterPage(page)
```

## Настройка админ-панели

```go
admin := formist.New().
    SetTitle("Моя Админ-панель").
    EnableAuth(true).
    EnableCORS(true, "http://localhost:3000").
    AddMiddleware(myMiddleware)
```

## Storage слой для хранения роутов

Библиотека поддерживает сохранение информации о роутах в базе данных для динамической навигации в UI.

### Подключение Storage

```go
import (
    "github.com/koteyye/go-formist"
    "github.com/koteyye/go-formist/storage/postgres"
)

// Создаем storage
storage, err := postgres.NewPostgresStorage(ctx, "postgres://user:pass@localhost/db")
if err != nil {
    log.Fatal(err)
}
defer storage.Close()

// Подключаем к админке
admin := formist.New().
    WithStorage(storage)
```

### PostgreSQL реализация

По умолчанию доступна реализация для PostgreSQL. При первом запуске автоматически создается таблица `formist_routes`.

Требования:
- PostgreSQL 12+
- Драйвер pgx/v5

### Создание собственной реализации

Вы можете создать свою реализацию интерфейса `storage.Storage`:

```go
type Storage interface {
    // SaveRoute сохраняет или обновляет роут
    SaveRoute(ctx context.Context, route *Route) error
    
    // GetRoutes возвращает все роуты для UI
    GetRoutes(ctx context.Context) ([]*Route, error)
    
    // DeleteRoute удаляет роут по ID
    DeleteRoute(ctx context.Context, id string) error
    
    // Close закрывает соединение
    Close() error
}
```

Пример реализации для MongoDB, Redis или любой другой БД можно найти в документации.

### API для работы с роутами

При подключенном Storage автоматически добавляются endpoints:

- `GET /api/routes` - получить все роуты из БД
- `POST /api/routes` - создать новый роут
- `PUT /api/routes/{id}` - обновить роут
- `DELETE /api/routes/delete?id={id}` - удалить роут

## API Endpoints

После запуска сервера доступны следующие endpoints:

- `GET /admin/config` - конфигурация админ-панели
- `GET /admin/forms/` - список форм
- `GET /admin/forms/{name}` - получение схемы формы
- `POST /admin/forms/{name}` - отправка данных формы
- `GET /admin/pages/{name}` - получение страницы

## Интеграция с фронтендом

Библиотека генерирует JSON Schema и UI Schema, которые можно использовать с любым фронтенд-фреймворком, поддерживающим JSON Schema Forms (например, React JSON Schema Form, Vue JSON Schema Form).

Пример ответа для формы:

```json
{
  "success": true,
  "data": {
    "schema": {
      "$schema": "https://json-schema.org/draft/2020-12/schema",
      "type": "object",
      "title": "Пользователь",
      "properties": {
        "name": {
          "type": "string",
          "title": "Имя"
        },
        "email": {
          "type": "string",
          "format": "email",
          "title": "Email"
        }
      },
      "required": ["name", "email"]
    },
    "uiSchema": {
      "ui:order": ["name", "email"],
      "email": {
        "ui:widget": "email"
      }
    }
  }
}
```

## Примеры

Больше примеров можно найти в папке [example/](example/).

### Быстрый запуск с Docker

Для быстрого локального запуска примера админки с PostgreSQL:

```bash
# Клонируем репозиторий
git clone https://github.com/koteyye/go-formist.git
cd go-formist

# Запускаем через docker-compose
docker-compose up -d

# Админка будет доступна на http://localhost:8080
# PostgreSQL на localhost:5432 (user: user, password: password)
```

Остановка:
```bash
docker-compose down

# Для удаления данных БД
docker-compose down -v
```

## Лицензия

MIT License. См. [LICENSE](LICENSE) для подробностей.
