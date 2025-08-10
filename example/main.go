package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/koteyye/go-formist"
	"github.com/koteyye/go-formist/types"
)

func main() {
	// Создаем новую админ-панель
	admin := formist.New()

	// Настраиваем админ-панель
	admin.SetTitle("Моя Админ-панель").
		EnableCORS(true, "http://localhost:3000").
		EnableAuth(false)

	// Создаем форму пользователя
	userForm := formist.NewForm("user", "Пользователь").
		WithDescription("Управление пользователями системы").
		AddTextField("name", "Имя").
		AddEmailField("email", "Email").
		AddPasswordField("password", "Пароль").
		AddSelectField("role", "Роль", []types.SelectOption{
			formist.SelectOption("admin", "Администратор"),
			formist.SelectOption("user", "Пользователь"),
			formist.SelectOption("moderator", "Модератор"),
		}).
		AddCheckboxField("active", "Активен").
		AddTextareaField("bio", "Биография").
		OnGet(func() (interface{}, error) {
			// Возвращаем тестовые данные
			return map[string]interface{}{
				"name":   "Иван Иванов",
				"email":  "ivan@example.com",
				"role":   "user",
				"active": true,
				"bio":    "Тестовый пользователь",
			}, nil
		}).
		OnPost(func(data map[string]interface{}) (interface{}, error) {
			// Обрабатываем данные формы
			fmt.Printf("Получены данные: %+v\n", data)
			return map[string]interface{}{
				"id":      123,
				"message": "Пользователь успешно создан",
			}, nil
		}).
		Build()

	// Создаем форму с таблицей
	ordersFormBuilder := formist.NewForm("orders", "Заказы").
		WithDescription("Управление заказами").
		AddTextField("search", "Поиск")

	// Добавляем таблицу как поле
	tableField := ordersFormBuilder.AddTableField("orders_table", "Список заказов").
		AddTextColumn("id", "ID").WithSortable().
		AddTextColumn("customer", "Клиент").WithFilterable().
		AddEmailColumn("email", "Email").
		AddNumberColumn("amount", "Сумма").WithSortable().
		AddSelectColumn("status", "Статус", []types.SelectOption{
			formist.SelectOption("pending", "В ожидании"),
			formist.SelectOption("processing", "Обрабатывается"),
			formist.SelectOption("completed", "Завершен"),
			formist.SelectOption("cancelled", "Отменен"),
		}).WithFilterable().
		AddDateColumn("created_at", "Дата создания").WithSortable().
		WithPagination(true).
		WithPageSize(20).
		WithSelectable(true).
		OnGet(func(page, limit int, filters map[string]interface{}) (types.TableData, error) {
			// Генерируем тестовые данные таблицы
			rows := []map[string]interface{}{
				{
					"id":         1,
					"customer":   "Петр Петров",
					"email":      "petr@example.com",
					"amount":     1500.50,
					"status":     "completed",
					"created_at": "2024-01-15",
				},
				{
					"id":         2,
					"customer":   "Анна Сидорова",
					"email":      "anna@example.com",
					"amount":     2300.00,
					"status":     "processing",
					"created_at": "2024-01-16",
				},
				{
					"id":         3,
					"customer":   "Михаил Иванов",
					"email":      "mikhail@example.com",
					"amount":     750.25,
					"status":     "pending",
					"created_at": "2024-01-17",
				},
			}

			return types.TableData{
				Columns: []types.TableColumn{
					{Key: "id", Title: "ID", Type: types.FieldTypeNumber, Sortable: true},
					{Key: "customer", Title: "Клиент", Type: types.FieldTypeText, Filterable: true},
					{Key: "email", Title: "Email", Type: types.FieldTypeEmail},
					{Key: "amount", Title: "Сумма", Type: types.FieldTypeNumber, Sortable: true},
					{Key: "status", Title: "Статус", Type: types.FieldTypeSelect, Filterable: true},
					{Key: "created_at", Title: "Дата создания", Type: types.FieldTypeDate, Sortable: true},
				},
				Rows:  rows,
				Total: 3,
				Page:  page,
				Limit: limit,
			}, nil
		})

	// Завершаем построение формы с таблицей
	ordersForm := tableField.Build(ordersFormBuilder).Build()

	// Создаем форму из структуры
	type Product struct {
		Name        string  `form:"name" label:"Название" required:"true"`
		Price       float64 `form:"price" label:"Цена" required:"true"`
		Description string  `form:"description" label:"Описание"`
		InStock     bool    `form:"in_stock" label:"В наличии"`
	}

	productForm := formist.FromStruct("products", "Товары", Product{}).
		OnPost(func(data map[string]interface{}) (interface{}, error) {
			fmt.Printf("Данные товара: %+v\n", data)
			return map[string]interface{}{
				"id":      456,
				"message": "Товар успешно создан",
			}, nil
		}).
		Build()

	// Создаем кастомную страницу
	dashboardPage := formist.NewPage("dashboard", "Панель управления").
		WithContent(`
			<h1>Добро пожаловать в админку!</h1>
			<p>Это кастомная HTML страница.</p>
			<ul>
				<li><a href="/admin/forms/user">Пользователи</a></li>
				<li><a href="/admin/forms/orders">Заказы</a></li>
				<li><a href="/admin/forms/products">Товары</a></li>
			</ul>
		`).
		Build()

	// Регистрируем формы и страницы
	admin.RegisterForm(userForm)
	admin.RegisterForm(ordersForm)
	admin.RegisterForm(productForm)
	admin.RegisterPage(dashboardPage)

	// Запускаем сервер
	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Админка доступна по адресу: http://localhost:8080/admin/config")
	fmt.Println("Формы:")
	fmt.Println("  - Пользователи: http://localhost:8080/admin/forms/user")
	fmt.Println("  - Заказы: http://localhost:8080/admin/forms/orders")
	fmt.Println("  - Товары: http://localhost:8080/admin/forms/products")
	fmt.Println("Страницы:")
	fmt.Println("  - Панель управления: http://localhost:8080/admin/pages/dashboard")

	log.Fatal(http.ListenAndServe(":8080", admin.Handler()))
}
