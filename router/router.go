package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/koteyye/go-formist/schema"
	"github.com/koteyye/go-formist/types"
)

// Router представляет HTTP роутер для админки
type Router struct {
	mux             *chi.Mux
	forms           map[string]*types.Form
	pages           map[string]*types.Page
	title           string
	authEnabled     bool
	corsEnabled     bool
	corsOrigins     []string
	middlewares     []types.MiddlewareFunc
	storageHandlers map[string]http.HandlerFunc
}

// NewRouter создает новый роутер
func NewRouter() *Router {
	r := &Router{
		mux:         chi.NewRouter(),
		forms:       make(map[string]*types.Form),
		pages:       make(map[string]*types.Page),
		title:       "Admin Panel",
		authEnabled: false,
		corsEnabled: false,
		corsOrigins: []string{"*"},
		middlewares: make([]types.MiddlewareFunc, 0),
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

// SetTitle устанавливает заголовок админки
func (r *Router) SetTitle(title string) {
	r.title = title
}

// EnableAuth включает авторизацию
func (r *Router) EnableAuth(enabled bool) {
	r.authEnabled = enabled
}

// EnableCORS включает CORS
func (r *Router) EnableCORS(enabled bool, origins ...string) {
	r.corsEnabled = enabled
	if len(origins) > 0 {
		r.corsOrigins = origins
	}
	// Пересоздаем mux с новыми настройками
	r.mux = chi.NewRouter()
	r.setupMiddleware()
	r.setupRoutes()
}

// AddMiddleware добавляет middleware
func (r *Router) AddMiddleware(middleware types.MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

// RegisterForm регистрирует форму
func (r *Router) RegisterForm(form *types.Form) {
	r.forms[form.Name] = form
}

// RegisterPage регистрирует страницу
func (r *Router) RegisterPage(page *types.Page) {
	r.pages[page.Name] = page
}

// Handler возвращает HTTP handler
func (r *Router) Handler() http.Handler {
	return r.mux
}

// SetStorageHandlers устанавливает обработчики для работы со storage
func (r *Router) SetStorageHandlers(handlers map[string]http.HandlerFunc) {
	r.storageHandlers = handlers
}

// setupMiddleware настраивает middleware
func (r *Router) setupMiddleware() {
	// Базовые middleware
	r.mux.Use(middleware.Logger)
	r.mux.Use(middleware.Recoverer)
	r.mux.Use(middleware.RequestID)

	// CORS
	if r.corsEnabled {
		r.mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   r.corsOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	// Кастомные middleware
	for _, mw := range r.middlewares {
		r.mux.Use(mw)
	}
}

// setupRoutes настраивает маршруты
func (r *Router) setupRoutes() {
	r.mux.Route("/admin", func(adminRouter chi.Router) {
		// Конфигурация админки
		adminRouter.Get("/config", r.handleConfig)

		// Формы
		adminRouter.Route("/forms", func(formsRouter chi.Router) {
			formsRouter.Get("/", r.handleFormsList)
			formsRouter.Get("/{name}", r.handleFormGet)
			formsRouter.Post("/{name}", r.handleFormPost)
		})

		// Страницы
		adminRouter.Route("/pages", func(pagesRouter chi.Router) {
			pagesRouter.Get("/{name}", r.handlePageGet)
		})

		// Авторизация (если включена)
		if r.authEnabled {
			adminRouter.Post("/login", r.handleLogin)
			adminRouter.Post("/logout", r.handleLogout)
		}
	})

	// API роуты (вне /admin для удобства)
	r.mux.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/routes", func(routesRouter chi.Router) {
			// GET /api/routes - получить все роуты
			routesRouter.Get("/", func(w http.ResponseWriter, req *http.Request) {
				if r.storageHandlers != nil {
					if handler, ok := r.storageHandlers["getRoutes"]; ok {
						handler(w, req)
						return
					}
				}
				http.Error(w, "Storage not configured", http.StatusNotImplemented)
			})

			// GET /api/routes/{id} - получить роут по ID
			routesRouter.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
				if r.storageHandlers != nil {
					if handler, ok := r.storageHandlers["getRoute"]; ok {
						handler(w, req)
						return
					}
				}
				http.Error(w, "Storage not configured", http.StatusNotImplemented)
			})

			// POST /api/routes - создать новый роут
			routesRouter.Post("/", func(w http.ResponseWriter, req *http.Request) {
				if r.storageHandlers != nil {
					if handler, ok := r.storageHandlers["createRoute"]; ok {
						handler(w, req)
						return
					}
				}
				http.Error(w, "Storage not configured", http.StatusNotImplemented)
			})

			// PUT /api/routes/{id} - обновить роут
			routesRouter.Put("/{id}", func(w http.ResponseWriter, req *http.Request) {
				if r.storageHandlers != nil {
					if handler, ok := r.storageHandlers["updateRoute"]; ok {
						handler(w, req)
						return
					}
				}
				http.Error(w, "Storage not configured", http.StatusNotImplemented)
			})

			// DELETE /api/routes/{id} - удалить роут
			routesRouter.Delete("/{id}", func(w http.ResponseWriter, req *http.Request) {
				if r.storageHandlers != nil {
					if handler, ok := r.storageHandlers["deleteRoute"]; ok {
						handler(w, req)
						return
					}
				}
				http.Error(w, "Storage not configured", http.StatusNotImplemented)
			})
		})
	})
}

// handleConfig обрабатывает запрос конфигурации
func (r *Router) handleConfig(w http.ResponseWriter, req *http.Request) {
	formsMap := make(map[string]string)
	for name, form := range r.forms {
		formsMap[name] = form.Title
	}

	pagesMap := make(map[string]string)
	for name, page := range r.pages {
		pagesMap[name] = page.Title
	}

	config := types.ConfigResponse{
		Title:       r.title,
		AuthEnabled: r.authEnabled,
		Forms:       formsMap,
		Pages:       pagesMap,
	}

	r.sendJSON(w, types.APIResponse{
		Success: true,
		Data:    config,
	})
}

// handleFormsList обрабатывает запрос списка форм
func (r *Router) handleFormsList(w http.ResponseWriter, req *http.Request) {
	formsMap := make(map[string]string)
	for name, form := range r.forms {
		formsMap[name] = form.Title
	}

	r.sendJSON(w, types.APIResponse{
		Success: true,
		Data:    formsMap,
	})
}

// handleFormGet обрабатывает GET запрос формы
func (r *Router) handleFormGet(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	form, exists := r.forms[name]
	if !exists {
		r.sendError(w, http.StatusNotFound, "Форма не найдена")
		return
	}

	// Генерируем схемы
	jsonSchema, err := schema.GenerateJSONSchema(form)
	if err != nil {
		r.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка генерации схемы: %v", err))
		return
	}

	uiSchema := schema.GenerateUISchema(form)

	response := types.FormResponse{
		Schema:   jsonSchema,
		UISchema: uiSchema,
	}

	// Если есть обработчик GET, получаем данные
	if form.OnGet != nil {
		data, err := form.OnGet()
		if err != nil {
			r.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка получения данных: %v", err))
			return
		}
		response.Data = data
	}

	r.sendJSON(w, types.APIResponse{
		Success: true,
		Data:    response,
	})
}

// handleFormPost обрабатывает POST запрос формы
func (r *Router) handleFormPost(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	form, exists := r.forms[name]
	if !exists {
		r.sendError(w, http.StatusNotFound, "Форма не найдена")
		return
	}

	if form.OnPost == nil {
		r.sendError(w, http.StatusMethodNotAllowed, "POST не поддерживается для этой формы")
		return
	}

	// Парсим данные
	var data map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		r.sendError(w, http.StatusBadRequest, "Некорректные данные JSON")
		return
	}

	// Валидируем данные
	if err := r.validateFormData(form, data); err != nil {
		r.sendError(w, http.StatusBadRequest, fmt.Sprintf("Ошибка валидации: %v", err))
		return
	}

	// Обрабатываем данные
	result, err := form.OnPost(data)
	if err != nil {
		r.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка обработки: %v", err))
		return
	}

	r.sendJSON(w, types.APIResponse{
		Success: true,
		Data:    result,
	})
}

// handlePageGet обрабатывает GET запрос страницы
func (r *Router) handlePageGet(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	page, exists := r.pages[name]
	if !exists {
		r.sendError(w, http.StatusNotFound, "Страница не найдена")
		return
	}

	// Если есть кастомный обработчик, используем его
	if page.Handler != nil {
		page.Handler(w, req)
		return
	}

	// Иначе возвращаем содержимое страницы
	r.sendJSON(w, types.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"title":   page.Title,
			"content": page.Content,
		},
	})
}

// handleLogin обрабатывает авторизацию
func (r *Router) handleLogin(w http.ResponseWriter, req *http.Request) {
	// TODO: Реализовать авторизацию
	r.sendJSON(w, types.APIResponse{
		Success: true,
		Message: "Авторизация успешна",
	})
}

// handleLogout обрабатывает выход
func (r *Router) handleLogout(w http.ResponseWriter, req *http.Request) {
	// TODO: Реализовать выход
	r.sendJSON(w, types.APIResponse{
		Success: true,
		Message: "Выход выполнен",
	})
}

// validateFormData валидирует данные формы
func (r *Router) validateFormData(form *types.Form, data map[string]interface{}) error {
	for _, field := range form.Fields {
		value, exists := data[field.Name]

		// Проверяем обязательные поля
		if field.Required && (!exists || isEmpty(value)) {
			return fmt.Errorf("поле '%s' обязательно для заполнения", field.Label)
		}

		// Если поле не обязательное и пустое, пропускаем валидацию
		if !exists || isEmpty(value) {
			continue
		}

		// Применяем правила валидации
		for _, rule := range field.Validation {
			if err := r.validateRule(value, rule); err != nil {
				return fmt.Errorf("поле '%s': %v", field.Label, err)
			}
		}
	}

	return nil
}

// validateRule применяет правило валидации
func (r *Router) validateRule(value interface{}, rule types.ValidationRule) error {
	switch rule.Type {
	case "email":
		return r.validateEmail(value, rule.Message)
	case "min":
		return r.validateMin(value, rule.Value, rule.Message)
	case "max":
		return r.validateMax(value, rule.Value, rule.Message)
	case "minLength":
		return r.validateMinLength(value, rule.Value, rule.Message)
	case "maxLength":
		return r.validateMaxLength(value, rule.Value, rule.Message)
	default:
		return nil
	}
}

// validateEmail валидирует email
func (r *Router) validateEmail(value interface{}, message string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("значение должно быть строкой")
	}

	if !strings.Contains(str, "@") || !strings.Contains(str, ".") {
		if message != "" {
			return fmt.Errorf("%s", message)
		}
		return fmt.Errorf("некорректный email адрес")
	}

	return nil
}

// validateMin валидирует минимальное значение
func (r *Router) validateMin(value interface{}, minValue interface{}, message string) error {
	num, err := toFloat64(value)
	if err != nil {
		return err
	}

	min, err := toFloat64(minValue)
	if err != nil {
		return err
	}

	if num < min {
		if message != "" {
			return fmt.Errorf("%s", message)
		}
		return fmt.Errorf("значение должно быть не менее %v", min)
	}

	return nil
}

// validateMax валидирует максимальное значение
func (r *Router) validateMax(value interface{}, maxValue interface{}, message string) error {
	num, err := toFloat64(value)
	if err != nil {
		return err
	}

	max, err := toFloat64(maxValue)
	if err != nil {
		return err
	}

	if num > max {
		if message != "" {
			return fmt.Errorf("%s", message)
		}
		return fmt.Errorf("значение должно быть не более %v", max)
	}

	return nil
}

// validateMinLength валидирует минимальную длину
func (r *Router) validateMinLength(value interface{}, minLength interface{}, message string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("значение должно быть строкой")
	}

	min, err := toInt(minLength)
	if err != nil {
		return err
	}

	if len(str) < min {
		if message != "" {
			return fmt.Errorf("%s", message)
		}
		return fmt.Errorf("длина должна быть не менее %d символов", min)
	}

	return nil
}

// validateMaxLength валидирует максимальную длину
func (r *Router) validateMaxLength(value interface{}, maxLength interface{}, message string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("значение должно быть строкой")
	}

	max, err := toInt(maxLength)
	if err != nil {
		return err
	}

	if len(str) > max {
		if message != "" {
			return fmt.Errorf("%s", message)
		}
		return fmt.Errorf("длина должна быть не более %d символов", max)
	}

	return nil
}

// sendJSON отправляет JSON ответ
func (r *Router) sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// sendError отправляет ошибку
func (r *Router) sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(types.APIResponse{
		Success: false,
		Error:   message,
	})
}

// isEmpty проверяет, является ли значение пустым
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []interface{}:
		return len(v) == 0
	default:
		return false
	}
}

// toFloat64 конвертирует значение в float64
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("не удается конвертировать %T в число", value)
	}
}

// toInt конвертирует значение в int
func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("не удается конвертировать %T в целое число", value)
	}
}
