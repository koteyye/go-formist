package form

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/koteyye/go-formist/types"
)

// FromStruct создает форму из Go структуры
func FromStruct(name, title string, structType interface{}) *FormBuilder {
	fb := NewForm(name, title)

	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fb
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Пропускаем неэкспортируемые поля
		if !field.IsExported() {
			continue
		}

		formField := createFieldFromStructField(field)
		if formField.Name != "" {
			fb.AddField(formField)
		}
	}

	return fb
}

// createFieldFromStructField создает поле формы из поля структуры
func createFieldFromStructField(field reflect.StructField) types.Field {
	formField := types.Field{
		Name:       getFieldName(field),
		Label:      getFieldLabel(field),
		Type:       getFieldType(field),
		Required:   getFieldRequired(field),
		Validation: make([]types.ValidationRule, 0),
	}

	// Добавляем валидацию для email полей
	if formField.Type == types.FieldTypeEmail {
		formField.Validation = append(formField.Validation, types.ValidationRule{
			Type:    "email",
			Message: "Введите корректный email",
		})
	}

	return formField
}

// getFieldName получает имя поля из тега form или имени поля
func getFieldName(field reflect.StructField) string {
	if name := field.Tag.Get("form"); name != "" {
		return name
	}
	return strings.ToLower(field.Name)
}

// getFieldLabel получает метку поля из тега label или имени поля
func getFieldLabel(field reflect.StructField) string {
	if label := field.Tag.Get("label"); label != "" {
		return label
	}
	return field.Name
}

// getFieldType определяет тип поля по типу Go и тегу type
func getFieldType(field reflect.StructField) types.FieldType {
	// Проверяем тег type
	if fieldType := field.Tag.Get("type"); fieldType != "" {
		switch fieldType {
		case "email":
			return types.FieldTypeEmail
		case "password":
			return types.FieldTypePassword
		case "textarea":
			return types.FieldTypeTextarea
		case "select":
			return types.FieldTypeSelect
		case "radio":
			return types.FieldTypeRadio
		case "checkbox":
			return types.FieldTypeCheckbox
		case "date":
			return types.FieldTypeDate
		case "time":
			return types.FieldTypeTime
		case "file":
			return types.FieldTypeFile
		case "hidden":
			return types.FieldTypeHidden
		case "number":
			return types.FieldTypeNumber
		}
	}

	// Определяем по типу Go
	switch field.Type.Kind() {
	case reflect.Bool:
		return types.FieldTypeCheckbox
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return types.FieldTypeNumber
	case reflect.String:
		// Проверяем имя поля для определения типа
		fieldName := strings.ToLower(field.Name)
		if strings.Contains(fieldName, "email") {
			return types.FieldTypeEmail
		}
		if strings.Contains(fieldName, "password") {
			return types.FieldTypePassword
		}
		return types.FieldTypeText
	default:
		return types.FieldTypeText
	}
}

// getFieldRequired проверяет, является ли поле обязательным
func getFieldRequired(field reflect.StructField) bool {
	required := field.Tag.Get("required")
	return required == "true" || required == "1"
}

// NewPage создает новую страницу
func NewPage(name, title string) *PageBuilder {
	return &PageBuilder{
		page: &types.Page{
			Name:  name,
			Title: title,
		},
	}
}

// PageBuilder представляет строитель страниц
type PageBuilder struct {
	page *types.Page
}

// WithContent устанавливает содержимое страницы
func (pb *PageBuilder) WithContent(content string) *PageBuilder {
	pb.page.Content = content
	return pb
}

// Build завершает построение страницы
func (pb *PageBuilder) Build() *types.Page {
	return pb.page
}

// ValidateField валидирует значение поля
func ValidateField(field *types.Field, value interface{}) error {
	// Проверка обязательного поля
	if field.Required && isEmpty(value) {
		return errors.New("поле обязательно для заполнения")
	}

	// Если поле пустое и не обязательное, пропускаем валидацию
	if isEmpty(value) {
		return nil
	}

	// Применяем правила валидации
	for _, rule := range field.Validation {
		if err := validateRule(value, rule); err != nil {
			return err
		}
	}

	return nil
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

// validateRule применяет правило валидации
func validateRule(value interface{}, rule types.ValidationRule) error {
	switch rule.Type {
	case "email":
		return validateEmail(value, rule.Message)
	case "min":
		return validateMin(value, rule.Value, rule.Message)
	case "max":
		return validateMax(value, rule.Value, rule.Message)
	case "minLength":
		return validateMinLength(value, rule.Value, rule.Message)
	case "maxLength":
		return validateMaxLength(value, rule.Value, rule.Message)
	case "pattern":
		return validatePattern(value, rule.Value, rule.Message)
	default:
		return nil
	}
}

// validateEmail валидирует email
func validateEmail(value interface{}, message string) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("значение должно быть строкой")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		if message != "" {
			return errors.New(message)
		}
		return errors.New("некорректный email адрес")
	}

	return nil
}

// validateMin валидирует минимальное значение
func validateMin(value interface{}, minValue interface{}, message string) error {
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
			return errors.New(message)
		}
		return fmt.Errorf("значение должно быть не менее %v", min)
	}

	return nil
}

// validateMax валидирует максимальное значение
func validateMax(value interface{}, maxValue interface{}, message string) error {
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
			return errors.New(message)
		}
		return fmt.Errorf("значение должно быть не более %v", max)
	}

	return nil
}

// validateMinLength валидирует минимальную длину строки
func validateMinLength(value interface{}, minLength interface{}, message string) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("значение должно быть строкой")
	}

	min, err := toInt(minLength)
	if err != nil {
		return err
	}

	if len(str) < min {
		if message != "" {
			return errors.New(message)
		}
		return fmt.Errorf("длина должна быть не менее %d символов", min)
	}

	return nil
}

// validateMaxLength валидирует максимальную длину строки
func validateMaxLength(value interface{}, maxLength interface{}, message string) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("значение должно быть строкой")
	}

	max, err := toInt(maxLength)
	if err != nil {
		return err
	}

	if len(str) > max {
		if message != "" {
			return errors.New(message)
		}
		return fmt.Errorf("длина должна быть не более %d символов", max)
	}

	return nil
}

// validatePattern валидирует по регулярному выражению
func validatePattern(value interface{}, pattern interface{}, message string) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("значение должно быть строкой")
	}

	patternStr, ok := pattern.(string)
	if !ok {
		return errors.New("паттерн должен быть строкой")
	}

	regex, err := regexp.Compile(patternStr)
	if err != nil {
		return fmt.Errorf("некорректное регулярное выражение: %v", err)
	}

	if !regex.MatchString(str) {
		if message != "" {
			return errors.New(message)
		}
		return errors.New("значение не соответствует требуемому формату")
	}

	return nil
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
