package types

import "net/http"

// FieldType представляет тип поля формы
type FieldType string

// Константы типов полей
const (
	FieldTypeText     FieldType = "text"
	FieldTypeEmail    FieldType = "email"
	FieldTypePassword FieldType = "password"
	FieldTypeNumber   FieldType = "number"
	FieldTypeTextarea FieldType = "textarea"
	FieldTypeSelect   FieldType = "select"
	FieldTypeRadio    FieldType = "radio"
	FieldTypeCheckbox FieldType = "checkbox"
	FieldTypeDate     FieldType = "date"
	FieldTypeTime     FieldType = "time"
	FieldTypeFile     FieldType = "file"
	FieldTypeHidden   FieldType = "hidden"
	FieldTypeTable    FieldType = "table"
)

// SelectOption представляет опцию для select/radio полей
type SelectOption struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled,omitempty"`
}

// ValidationRule представляет правило валидации
type ValidationRule struct {
	Type    string      `json:"type"`
	Value   interface{} `json:"value,omitempty"`
	Message string      `json:"message"`
}

// TableColumn представляет колонку таблицы
type TableColumn struct {
	Key        string         `json:"key"`
	Title      string         `json:"title"`
	Type       FieldType      `json:"type"`
	Sortable   bool           `json:"sortable,omitempty"`
	Filterable bool           `json:"filterable,omitempty"`
	Width      string         `json:"width,omitempty"`
	Align      string         `json:"align,omitempty"`
	Options    []SelectOption `json:"options,omitempty"`
	Multiple   bool           `json:"multiple,omitempty"`
}

// TableData представляет данные таблицы
type TableData struct {
	Columns []TableColumn            `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
	Total   int                      `json:"total"`
	Page    int                      `json:"page"`
	Limit   int                      `json:"limit"`
}

// TableConfig представляет конфигурацию таблицы
type TableConfig struct {
	Columns    []TableColumn `json:"columns"`
	Pagination bool          `json:"pagination"`
	PageSize   int           `json:"pageSize"`
	Sortable   bool          `json:"sortable"`
	Filterable bool          `json:"filterable"`
	Selectable bool          `json:"selectable"`
	Editable   bool          `json:"editable"`
	OnGet      TableHandler  `json:"-"`
}

// Field представляет поле формы
type Field struct {
	Name         string                 `json:"name"`
	Type         FieldType              `json:"type"`
	Label        string                 `json:"label"`
	Required     bool                   `json:"required"`
	Placeholder  string                 `json:"placeholder,omitempty"`
	DefaultValue interface{}            `json:"defaultValue,omitempty"`
	Options      []SelectOption         `json:"options,omitempty"`
	Multiple     bool                   `json:"multiple,omitempty"`
	Validation   []ValidationRule       `json:"validation,omitempty"`
	Group        string                 `json:"group,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Disabled     bool                   `json:"disabled,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
	TableConfig  *TableConfig           `json:"tableConfig,omitempty"`
}

// FieldGroup представляет группу полей
type FieldGroup struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Fields      []string `json:"fields"`
}

// Form представляет форму
type Form struct {
	Name        string       `json:"name"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Fields      []Field      `json:"fields"`
	Groups      []FieldGroup `json:"groups,omitempty"`
	OnPost      FormHandler  `json:"-"`
	OnGet       GetHandler   `json:"-"`
}

// Page представляет кастомную страницу
type Page struct {
	Name    string             `json:"name"`
	Title   string             `json:"title"`
	Content string             `json:"content,omitempty"`
	Handler http.HandlerFunc   `json:"-"`
}

// Обработчики
type FormHandler func(data map[string]interface{}) (interface{}, error)
type GetHandler func() (interface{}, error)
type TableHandler func(page, limit int, filters map[string]interface{}) (TableData, error)
type MiddlewareFunc func(http.Handler) http.Handler

// API Response структуры
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ConfigResponse struct {
	Title       string            `json:"title"`
	AuthEnabled bool              `json:"authEnabled"`
	Forms       map[string]string `json:"forms"`
	Pages       map[string]string `json:"pages"`
}

type FormResponse struct {
	Schema   interface{} `json:"schema"`
	UISchema interface{} `json:"uiSchema"`
	Data     interface{} `json:"data,omitempty"`
}
