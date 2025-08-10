package schema

import (
	"fmt"

	"github.com/koteyye/go-formist/types"
)

// JSONSchema представляет JSON Schema v7
type JSONSchema struct {
	Schema      string                 `json:"$schema"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Definitions map[string]interface{} `json:"definitions,omitempty"`
}

// UISchema представляет UI Schema для рендеринга
type UISchema struct {
	UIOrder   []string               `json:"ui:order,omitempty"`
	UIOptions map[string]interface{} `json:"ui:options,omitempty"`
	Fields    map[string]interface{} `json:",inline"`
}

// GenerateJSONSchema генерирует JSON Schema из формы
func GenerateJSONSchema(form *types.Form) (*JSONSchema, error) {
	schema := &JSONSchema{
		Schema:      "https://json-schema.org/draft/2020-12/schema",
		Type:        "object",
		Title:       form.Title,
		Description: form.Description,
		Properties:  make(map[string]interface{}),
		Required:    make([]string, 0),
		Definitions: make(map[string]interface{}),
	}

	// Обрабатываем поля формы
	for _, field := range form.Fields {
		// Пропускаем скрытые поля в схеме
		if field.Type == types.FieldTypeHidden {
			continue
		}

		fieldSchema, err := generateFieldSchema(&field)
		if err != nil {
			return nil, fmt.Errorf("ошибка генерации схемы для поля %s: %w", field.Name, err)
		}

		schema.Properties[field.Name] = fieldSchema

		// Добавляем в обязательные поля
		if field.Required {
			schema.Required = append(schema.Required, field.Name)
		}
	}

	return schema, nil
}

// GenerateUISchema генерирует UI Schema из формы
func GenerateUISchema(form *types.Form) map[string]interface{} {
	uiSchema := make(map[string]interface{})

	// Порядок полей
	order := make([]string, 0)
	for _, field := range form.Fields {
		if field.Type != types.FieldTypeHidden {
			order = append(order, field.Name)
		}
	}
	uiSchema["ui:order"] = order

	// Настройки для каждого поля
	for _, field := range form.Fields {
		fieldUI := generateFieldUISchema(&field)
		if len(fieldUI) > 0 {
			uiSchema[field.Name] = fieldUI
		}
	}

	// Группы полей
	if len(form.Groups) > 0 {
		groups := make([]map[string]interface{}, 0)
		for _, group := range form.Groups {
			groupUI := map[string]interface{}{
				"ui:title":       group.Title,
				"ui:description": group.Description,
				"ui:fields":      group.Fields,
			}
			groups = append(groups, groupUI)
		}
		uiSchema["ui:groups"] = groups
	}

	return uiSchema
}

// generateFieldSchema генерирует схему для отдельного поля
func generateFieldSchema(field *types.Field) (map[string]interface{}, error) {
	fieldSchema := map[string]interface{}{
		"title":       field.Label,
		"description": field.Description,
	}

	// Устанавливаем тип и формат в зависимости от типа поля
	switch field.Type {
	case types.FieldTypeText:
		fieldSchema["type"] = "string"
		if field.Placeholder != "" {
			fieldSchema["examples"] = []string{field.Placeholder}
		}

	case types.FieldTypeEmail:
		fieldSchema["type"] = "string"
		fieldSchema["format"] = "email"

	case types.FieldTypePassword:
		fieldSchema["type"] = "string"
		fieldSchema["format"] = "password"

	case types.FieldTypeNumber:
		fieldSchema["type"] = "number"

	case types.FieldTypeTextarea:
		fieldSchema["type"] = "string"

	case types.FieldTypeDate:
		fieldSchema["type"] = "string"
		fieldSchema["format"] = "date"

	case types.FieldTypeTime:
		fieldSchema["type"] = "string"
		fieldSchema["format"] = "time"

	case types.FieldTypeFile:
		fieldSchema["type"] = "string"
		fieldSchema["format"] = "data-url"

	case types.FieldTypeCheckbox:
		fieldSchema["type"] = "boolean"

	case types.FieldTypeRadio, types.FieldTypeSelect:
		if field.Multiple {
			fieldSchema["type"] = "array"
			fieldSchema["items"] = map[string]interface{}{
				"type": "string",
				"enum": getOptionValues(field.Options),
			}
			fieldSchema["uniqueItems"] = true
		} else {
			fieldSchema["type"] = "string"
			fieldSchema["enum"] = getOptionValues(field.Options)
		}

	case types.FieldTypeTable:
		if field.TableConfig != nil {
			tableSchema := generateTableSchema(field.TableConfig)
			fieldSchema = tableSchema
		} else {
			fieldSchema["type"] = "object"
		}

	default:
		fieldSchema["type"] = "string"
	}

	// Добавляем значение по умолчанию
	if field.DefaultValue != nil {
		fieldSchema["default"] = field.DefaultValue
	}

	// Добавляем правила валидации
	for _, rule := range field.Validation {
		switch rule.Type {
		case "min":
			if num, ok := rule.Value.(float64); ok {
				fieldSchema["minimum"] = num
			}
		case "max":
			if num, ok := rule.Value.(float64); ok {
				fieldSchema["maximum"] = num
			}
		case "minLength":
			if num, ok := rule.Value.(float64); ok {
				fieldSchema["minLength"] = int(num)
			}
		case "maxLength":
			if num, ok := rule.Value.(float64); ok {
				fieldSchema["maxLength"] = int(num)
			}
		case "pattern":
			if pattern, ok := rule.Value.(string); ok {
				fieldSchema["pattern"] = pattern
			}
		}
	}

	return fieldSchema, nil
}

// generateFieldUISchema генерирует UI схему для отдельного поля
func generateFieldUISchema(field *types.Field) map[string]interface{} {
	uiSchema := make(map[string]interface{})

	// Настройки виджета в зависимости от типа поля
	switch field.Type {
	case types.FieldTypePassword:
		uiSchema["ui:widget"] = "password"

	case types.FieldTypeTextarea:
		uiSchema["ui:widget"] = "textarea"
		uiSchema["ui:options"] = map[string]interface{}{
			"rows": 4,
		}

	case types.FieldTypeFile:
		uiSchema["ui:widget"] = "file"

	case types.FieldTypeCheckbox:
		uiSchema["ui:widget"] = "checkbox"

	case types.FieldTypeRadio:
		uiSchema["ui:widget"] = "radio"
		if len(field.Options) > 0 {
			uiSchema["ui:options"] = map[string]interface{}{
				"enumOptions": convertOptionsToEnumOptions(field.Options),
			}
		}

	case types.FieldTypeSelect:
		if field.Multiple {
			uiSchema["ui:widget"] = "checkboxes"
		} else {
			uiSchema["ui:widget"] = "select"
		}
		if len(field.Options) > 0 {
			uiSchema["ui:options"] = map[string]interface{}{
				"enumOptions": convertOptionsToEnumOptions(field.Options),
			}
		}

	case types.FieldTypeTable:
		uiSchema["ui:widget"] = "table"
		if field.TableConfig != nil {
			uiSchema["ui:options"] = generateTableUIOptions(field.TableConfig)
		}

	case types.FieldTypeHidden:
		uiSchema["ui:widget"] = "hidden"
	}

	// Placeholder
	if field.Placeholder != "" {
		uiSchema["ui:placeholder"] = field.Placeholder
	}

	// Disabled
	if field.Disabled {
		uiSchema["ui:disabled"] = true
	}

	// Группа
	if field.Group != "" {
		uiSchema["ui:group"] = field.Group
	}

	// Дополнительные настройки из Config
	if len(field.Config) > 0 {
		if uiOptions, exists := uiSchema["ui:options"]; exists {
			if optionsMap, ok := uiOptions.(map[string]interface{}); ok {
				for key, value := range field.Config {
					optionsMap[key] = value
				}
			}
		} else {
			uiSchema["ui:options"] = field.Config
		}
	}

	return uiSchema
}

// generateTableSchema генерирует схему для таблицы
func generateTableSchema(config *types.TableConfig) map[string]interface{} {
	schema := map[string]interface{}{
		"type":  "object",
		"title": "Таблица",
		"properties": map[string]interface{}{
			"columns": map[string]interface{}{
				"type":  "array",
				"items": generateTableColumnSchema(),
			},
			"rows": map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{
					"type": "object",
				},
			},
			"total": map[string]interface{}{
				"type": "integer",
			},
			"page": map[string]interface{}{
				"type": "integer",
			},
			"limit": map[string]interface{}{
				"type": "integer",
			},
		},
	}

	return schema
}

// generateTableColumnSchema генерирует схему для колонки таблицы
func generateTableColumnSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"key": map[string]interface{}{
				"type": "string",
			},
			"title": map[string]interface{}{
				"type": "string",
			},
			"type": map[string]interface{}{
				"type": "string",
			},
			"sortable": map[string]interface{}{
				"type": "boolean",
			},
			"filterable": map[string]interface{}{
				"type": "boolean",
			},
			"width": map[string]interface{}{
				"type": "string",
			},
			"align": map[string]interface{}{
				"type": "string",
			},
		},
	}
}

// generateTableUIOptions генерирует UI опции для таблицы
func generateTableUIOptions(config *types.TableConfig) map[string]interface{} {
	options := map[string]interface{}{
		"pagination":  config.Pagination,
		"pageSize":    config.PageSize,
		"sortable":    config.Sortable,
		"filterable":  config.Filterable,
		"selectable":  config.Selectable,
		"editable":    config.Editable,
		"columns":     config.Columns,
	}

	return options
}

// getOptionValues извлекает значения из опций
func getOptionValues(options []types.SelectOption) []string {
	values := make([]string, len(options))
	for i, option := range options {
		values[i] = option.Value
	}
	return values
}

// convertOptionsToEnumOptions конвертирует опции в формат enumOptions
func convertOptionsToEnumOptions(options []types.SelectOption) []map[string]interface{} {
	enumOptions := make([]map[string]interface{}, len(options))
	for i, option := range options {
		enumOptions[i] = map[string]interface{}{
			"value": option.Value,
			"label": option.Label,
		}
	}
	return enumOptions
}
