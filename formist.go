package formist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koteyye/go-formist/form"
	"github.com/koteyye/go-formist/router"
	"github.com/koteyye/go-formist/storage"
	"github.com/koteyye/go-formist/types"
)

// Admin представляет основной объект админ-панели с поддержкой storage
type Admin struct {
	router  *router.Router
	storage storage.Storage
}

// New создает новую админ-панель
func New() *Admin {
	return &Admin{
		router: router.NewRouter(),
	}
}

// WithStorage подключает storage для сохранения роутов
func (a *Admin) WithStorage(s storage.Storage) *Admin {
	a.storage = s
	return a
}

// SetTitle устанавливает заголовок админ-панели
func (a *Admin) SetTitle(title string) *Admin {
	a.router.SetTitle(title)
	return a
}

// EnableAuth включает авторизацию
func (a *Admin) EnableAuth(enabled bool) *Admin {
	a.router.EnableAuth(enabled)
	return a
}

// EnableCORS включает CORS
func (a *Admin) EnableCORS(enabled bool, origins ...string) *Admin {
	a.router.EnableCORS(enabled, origins...)
	return a
}

// AddMiddleware добавляет middleware
func (a *Admin) AddMiddleware(middleware types.MiddlewareFunc) *Admin {
	a.router.AddMiddleware(middleware)
	return a
}

// RegisterForm регистрирует форму и сохраняет роут в storage
func (a *Admin) RegisterForm(form *types.Form) *Admin {
	a.router.RegisterForm(form)
	
	// Сохраняем роут в storage если он подключен
	if a.storage != nil {
		route := &storage.Route{
			Name:  form.Name,
			Path:  fmt.Sprintf("/admin/forms/%s", form.Name),
			Title: form.Title,
			Type:  "form",
		}
		
		if form.Description != "" {
			route.Description = form.Description
		}
		
		// Игнорируем ошибку, чтобы не ломать работу если storage недоступен
		_ = a.storage.SaveRoute(context.Background(), route)
	}
	
	return a
}

// RegisterPage регистрирует страницу и сохраняет роут в storage
func (a *Admin) RegisterPage(page *types.Page) *Admin {
	a.router.RegisterPage(page)
	
	// Сохраняем роут в storage если он подключен
	if a.storage != nil {
		route := &storage.Route{
			Name:  page.Name,
			Path:  fmt.Sprintf("/admin/pages/%s", page.Name),
			Title: page.Title,
			Type:  "page",
		}
		
		// Игнорируем ошибку, чтобы не ломать работу если storage недоступен
		_ = a.storage.SaveRoute(context.Background(), route)
	}
	
	return a
}

// GetRoutes возвращает все роуты из storage
func (a *Admin) GetRoutes(ctx context.Context) ([]*storage.Route, error) {
	if a.storage == nil {
		return nil, fmt.Errorf("storage не подключен")
	}
	
	return a.storage.GetRoutes(ctx)
}

// DeleteRoute удаляет роут из storage
func (a *Admin) DeleteRoute(ctx context.Context, id string) error {
	if a.storage == nil {
		return fmt.Errorf("storage не подключен")
	}
	
	return a.storage.DeleteRoute(ctx, id)
}

// Handler возвращает HTTP handler для использования с любым HTTP сервером
func (a *Admin) Handler() http.Handler {
	// Добавляем эндпоинты для работы с роутами через storage
	if a.storage != nil {
		// Обертка над основным handler для добавления storage endpoints
		mux := http.NewServeMux()
		
		// Основные роуты админки
		mux.Handle("/admin/", a.router.Handler())
		
		// API для работы с роутами
		mux.HandleFunc("/api/routes", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				a.handleGetRoutes(w, r)
			case http.MethodDelete:
				a.handleDeleteRoute(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})
		
		return mux
	}
	
	return a.router.Handler()
}

// handleGetRoutes обрабатывает получение списка роутов
func (a *Admin) handleGetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := a.GetRoutes(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"routes":  routes,
	})
}

// handleDeleteRoute обрабатывает удаление роута
func (a *Admin) handleDeleteRoute(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	
	if err := a.DeleteRoute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// ListenAndServe запускает HTTP сервер на указанном адресе
func (a *Admin) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, a.Handler())
}

// Экспортируем основные функции для удобства использования

// NewForm создает новую форму
func NewForm(name, title string) *form.FormBuilder {
	return form.NewForm(name, title)
}

// NewPage создает новую страницу
func NewPage(name, title string) *form.PageBuilder {
	return form.NewPage(name, title)
}

// FromStruct создает форму из Go структуры
func FromStruct(name, title string, structType interface{}) *form.FormBuilder {
	return form.FromStruct(name, title, structType)
}

// SelectOption создает опцию для select/radio полей
func SelectOption(value, label string) types.SelectOption {
	return types.SelectOption{
		Value: value,
		Label: label,
	}
}

// ValidationRule создает правило валидации
func ValidationRule(ruleType string, value interface{}, message string) types.ValidationRule {
	return types.ValidationRule{
		Type:    ruleType,
		Value:   value,
		Message: message,
	}
}