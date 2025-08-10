package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/koteyye/go-formist"
	"github.com/koteyye/go-formist/storage/postgres"
	"github.com/koteyye/go-formist/types"
)

func main() {
	ctx := context.Background()

	// Получаем DSN из переменной окружения или используем дефолтный
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/formist_db?sslmode=disable"
	}

	// Создаем storage
	storage, err := postgres.NewPostgresStorage(ctx, dsn)
	if err != nil {
		log.Printf("Внимание: не удалось подключить storage: %v", err)
		log.Println("Продолжаем работу без сохранения роутов в БД")
	}

	// Создаем админ-панель
	admin := formist.New().
		SetTitle("Админка с хранилищем роутов").
		EnableCORS(true, "http://localhost:3000", "http://localhost:5173")

	// Подключаем storage если он доступен
	if storage != nil {
		admin.WithStorage(storage)
		defer storage.Close()
	}

	// Создаем форму пользователя
	userForm := formist.NewForm("users", "Пользователи").
		AddField(types.Field{
			Name:     "name",
			Type:     types.FieldTypeText,
			Label:    "Имя",
			Required: true,
		}).
		AddField(types.Field{
			Name:     "email",
			Type:     types.FieldTypeEmail,
			Label:    "Email",
			Required: true,
			Validation: []types.ValidationRule{
				{Type: "email", Message: "Введите корректный email"},
			},
		}).
		AddField(types.Field{
			Name:     "password",
			Type:     types.FieldTypePassword,
			Label:    "Пароль",
			Required: true,
		}).
		AddSelectField("role", "Роль", []types.SelectOption{
			formist.SelectOption("admin", "Администратор"),
			formist.SelectOption("user", "Пользователь"),
			formist.SelectOption("moderator", "Модератор"),
		}).
		AddCheckboxField("active", "Активен").
		OnPost(func(data map[string]interface{}) (interface{}, error) {
			// Здесь логика сохранения пользователя
			fmt.Printf("Создание пользователя: %+v\n", data)
			return map[string]string{
				"message": "Пользователь успешно создан",
				"id":      "user_123",
			}, nil
		}).
		Build()

	// Создаем форму настроек
	settingsForm := formist.NewForm("settings", "Настройки").
		AddTextField("site_name", "Название сайта").
		AddTextareaField("site_description", "Описание сайта").
		AddNumberField("items_per_page", "Элементов на странице").
		AddCheckboxField("maintenance_mode", "Режим обслуживания").
		OnGet(func() (interface{}, error) {
			// Загружаем текущие настройки
			return map[string]interface{}{
				"site_name":        "Мой сайт",
				"site_description": "Описание моего сайта",
				"items_per_page":   20,
				"maintenance_mode": false,
			}, nil
		}).
		OnPost(func(data map[string]interface{}) (interface{}, error) {
			fmt.Printf("Сохранение настроек: %+v\n", data)
			return map[string]string{"message": "Настройки сохранены"}, nil
		}).
		Build()

	// Создаем полную форму с использованием всех типов полей
	completeForm := formist.NewForm("complete_example", "Полный пример всех типов полей").
		// 1. Text field
		AddField(types.Field{
			Name:        "username",
			Type:        types.FieldTypeText,
			Label:       "Имя пользователя",
			Required:    true,
			Placeholder: "Введите имя пользователя",
		}).
		// 2. Email field
		AddField(types.Field{
			Name:     "user_email",
			Type:     types.FieldTypeEmail,
			Label:    "Email адрес",
			Required: true,
			Validation: []types.ValidationRule{
				{Type: "email", Message: "Введите корректный email"},
			},
		}).
		// 3. Password field
		AddField(types.Field{
			Name:        "user_password",
			Type:        types.FieldTypePassword,
			Label:       "Пароль",
			Required:    true,
			Placeholder: "Введите пароль",
		}).
		// 4. Number field
		AddField(types.Field{
			Name:         "age",
			Type:         types.FieldTypeNumber,
			Label:        "Возраст",
			Required:     false,
			DefaultValue: 18,
			Validation: []types.ValidationRule{
				{Type: "min", Value: 0, Message: "Возраст не может быть отрицательным"},
				{Type: "max", Value: 120, Message: "Возраст не может быть больше 120"},
			},
		}).
		// 5. Textarea field
		AddField(types.Field{
			Name:        "bio",
			Type:        types.FieldTypeTextarea,
			Label:       "Биография",
			Placeholder: "Расскажите о себе...",
			Config: map[string]interface{}{
				"rows": 4,
			},
		}).
		// 6. Select field
		AddField(types.Field{
			Name:     "country",
			Type:     types.FieldTypeSelect,
			Label:    "Страна",
			Required: true,
			Options: []types.SelectOption{
				{Value: "ru", Label: "Россия"},
				{Value: "us", Label: "США"},
				{Value: "de", Label: "Германия"},
				{Value: "fr", Label: "Франция"},
			},
		}).
		// 7. Radio field
		AddField(types.Field{
			Name:     "gender",
			Type:     types.FieldTypeRadio,
			Label:    "Пол",
			Required: true,
			Options: []types.SelectOption{
				{Value: "male", Label: "Мужской"},
				{Value: "female", Label: "Женский"},
				{Value: "other", Label: "Другой"},
			},
		}).
		// 8. Checkbox field
		AddField(types.Field{
			Name:         "newsletter",
			Type:         types.FieldTypeCheckbox,
			Label:        "Подписаться на рассылку",
			DefaultValue: false,
		}).
		// 9. Date field
		AddField(types.Field{
			Name:     "birth_date",
			Type:     types.FieldTypeDate,
			Label:    "Дата рождения",
			Required: false,
		}).
		// 10. Time field
		AddField(types.Field{
			Name:  "preferred_time",
			Type:  types.FieldTypeTime,
			Label: "Предпочитаемое время звонка",
		}).
		// 11. File field
		AddField(types.Field{
			Name:  "avatar",
			Type:  types.FieldTypeFile,
			Label: "Аватар",
			Config: map[string]interface{}{
				"accept": "image/*",
			},
		}).
		// 12. Hidden field
		AddField(types.Field{
			Name:         "user_id",
			Type:         types.FieldTypeHidden,
			DefaultValue: "hidden_user_123",
		}).
		// 13. Table field
		AddField(types.Field{
			Name:  "skills_table",
			Type:  types.FieldTypeTable,
			Label: "Навыки и опыт",
			TableConfig: &types.TableConfig{
				Columns: []types.TableColumn{
					{
						Key:   "skill",
						Title: "Навык",
						Type:  types.FieldTypeText,
						Width: "40%",
					},
					{
						Key:   "level",
						Title: "Уровень",
						Type:  types.FieldTypeSelect,
						Width: "30%",
						Options: []types.SelectOption{
							{Value: "beginner", Label: "Начинающий"},
							{Value: "intermediate", Label: "Средний"},
							{Value: "advanced", Label: "Продвинутый"},
							{Value: "expert", Label: "Эксперт"},
						},
					},
					{
						Key:   "years",
						Title: "Лет опыта",
						Type:  types.FieldTypeNumber,
						Width: "30%",
					},
				},
				Pagination: false,
				Sortable:   true,
				Editable:   true,
				OnGet: func(page, limit int, filters map[string]interface{}) (types.TableData, error) {
					return types.TableData{
						Columns: []types.TableColumn{
							{Key: "skill", Title: "Навык", Type: types.FieldTypeText},
							{Key: "level", Title: "Уровень", Type: types.FieldTypeSelect},
							{Key: "years", Title: "Лет опыта", Type: types.FieldTypeNumber},
						},
						Rows: []map[string]interface{}{
							{"skill": "Go", "level": "expert", "years": 5},
							{"skill": "JavaScript", "level": "advanced", "years": 3},
							{"skill": "Python", "level": "intermediate", "years": 2},
						},
						Total: 3,
						Page:  1,
						Limit: 10,
					}, nil
				},
			},
		}).
		OnGet(func() (interface{}, error) {
			// Возвращаем предзаполненные данные для демонстрации
			return map[string]interface{}{
				"username":       "demo_user",
				"user_email":     "demo@example.com",
				"age":            25,
				"bio":            "Это демонстрационная биография пользователя",
				"country":        "ru",
				"gender":         "male",
				"newsletter":     true,
				"birth_date":     "1999-01-15",
				"preferred_time": "14:30",
				"user_id":        "hidden_user_123",
			}, nil
		}).
		OnPost(func(data map[string]interface{}) (interface{}, error) {
			fmt.Printf("Полная форма отправлена: %+v\n", data)
			return map[string]interface{}{
				"message": "Форма успешно обработана!",
				"data":    data,
			}, nil
		}).
		Build()

	// Создаем кастомную страницу
	dashboardPage := formist.NewPage("dashboard", "Панель управления").
		WithContent(`
			<div class="dashboard">
				<h1>Добро пожаловать в админку!</h1>
				<p>Это кастомная страница с произвольным HTML контентом.</p>
				<div class="stats">
					<div class="stat-card">
						<h3>Пользователей</h3>
						<p class="stat-value">1,234</p>
					</div>
					<div class="stat-card">
						<h3>Заказов</h3>
						<p class="stat-value">567</p>
					</div>
				</div>
			</div>
		`).
		Build()

	// Регистрируем формы и страницы
	admin.RegisterForm(completeForm)
	admin.RegisterForm(userForm)
	admin.RegisterForm(settingsForm)
	admin.RegisterPage(dashboardPage)

	// Создаем главную страницу с документацией
	docsPage := formist.NewPage("docs", "Документация").
		WithContent(`
		<html>
		<head>
			<title>Go-Formist с Storage</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; }
				.endpoint { background: #f4f4f4; padding: 10px; margin: 10px 0; }
				code { background: #e0e0e0; padding: 2px 4px; }
			</style>
		</head>
		<body>
			<h1>Go-Formist с поддержкой Storage</h1>
			<h2>Доступные endpoints:</h2>
			
			<div class="endpoint">
				<strong>GET /admin/config</strong> - Конфигурация админки
			</div>
			
			<div class="endpoint">
				<strong>GET /admin/forms/</strong> - Список форм
			</div>
			
			<div class="endpoint">
				<strong>GET /admin/forms/{name}</strong> - Схема формы
			</div>
			
			<div class="endpoint">
				<strong>POST /admin/forms/{name}</strong> - Отправка данных формы
			</div>
			
			<div class="endpoint">
				<strong>GET /api/routes</strong> - Получить все роуты из БД
			</div>
			
			<div class="endpoint">
				<strong>GET /api/routes/{id}</strong> - Получить роут по ID
			</div>
			
			<div class="endpoint">
				<strong>POST /api/routes</strong> - Создать новый роут
			</div>
			
			<div class="endpoint">
				<strong>PUT /api/routes/{id}</strong> - Обновить роут
			</div>
			
			<div class="endpoint">
				<strong>DELETE /api/routes/{id}</strong> - Удалить роут
			</div>
			
			<h2>Пример использования с React:</h2>
			<pre><code>
// Получение роутов для навигации
fetch('http://localhost:8080/api/routes')
  .then(res => res.json())
  .then(data => {
    if (data.success) {
      const routes = data.routes;
      // Построение навигации на основе роутов
      routes.forEach(route => {
        console.log(route.title, route.path, route.type);
      });
    }
  });
			</code></pre>
		</body>
		</html>
		`).Build()

	// Регистрируем главную страницу
	admin.RegisterPage(docsPage)

	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Админка доступна на http://localhost:8080/admin/")
	fmt.Println("API роутов: http://localhost:8080/api/routes")

	log.Fatal(http.ListenAndServe(":8080", admin.Handler()))
}
