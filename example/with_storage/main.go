package main

import (
	"context"
	"encoding/json"
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
	admin.RegisterForm(userForm)
	admin.RegisterForm(settingsForm)
	admin.RegisterPage(dashboardPage)

	// Создаем дополнительный роутер для API роутов
	mux := http.NewServeMux()
	
	// Основной handler админки
	mux.Handle("/admin/", admin.Handler())
	
	// API endpoint для получения роутов (для UI)
	mux.HandleFunc("/api/routes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		routes, err := admin.GetRoutes(r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"routes":  routes,
		})
	})

	// API endpoint для удаления роута
	mux.HandleFunc("/api/routes/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		if err := admin.DeleteRoute(r.Context(), id); err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Роут успешно удален",
		})
	})

	// Главная страница с документацией
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
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
				<strong>DELETE /api/routes/delete?id={id}</strong> - Удалить роут
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
		`)
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Админка доступна на http://localhost:8080/admin/")
	fmt.Println("API роутов: http://localhost:8080/api/routes")
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}