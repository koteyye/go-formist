package form

import (
	"github.com/koteyye/go-formist/types"
)

// FormBuilder представляет строитель форм
type FormBuilder struct {
	form *types.Form
}

// TableFieldBuilder представляет строитель таблицы как поля формы
type TableFieldBuilder struct {
	field *types.Field
}

// NewForm создает новую форму
func NewForm(name, title string) *FormBuilder {
	return &FormBuilder{
		form: &types.Form{
			Name:   name,
			Title:  title,
			Fields: make([]types.Field, 0),
			Groups: make([]types.FieldGroup, 0),
		},
	}
}

// WithDescription устанавливает описание формы
func (fb *FormBuilder) WithDescription(description string) *FormBuilder {
	fb.form.Description = description
	return fb
}

// AddField добавляет поле в форму
func (fb *FormBuilder) AddField(field types.Field) *FormBuilder {
	fb.form.Fields = append(fb.form.Fields, field)
	return fb
}

// AddTextField добавляет текстовое поле
func (fb *FormBuilder) AddTextField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypeText,
		Label: label,
	}
	return fb.AddField(field)
}

// AddEmailField добавляет поле email
func (fb *FormBuilder) AddEmailField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypeEmail,
		Label: label,
		Validation: []types.ValidationRule{
			{Type: "email", Message: "Введите корректный email"},
		},
	}
	return fb.AddField(field)
}

// AddPasswordField добавляет поле пароля
func (fb *FormBuilder) AddPasswordField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypePassword,
		Label: label,
	}
	return fb.AddField(field)
}

// AddNumberField добавляет числовое поле
func (fb *FormBuilder) AddNumberField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypeNumber,
		Label: label,
	}
	return fb.AddField(field)
}

// AddSelectField добавляет поле выбора
func (fb *FormBuilder) AddSelectField(name, label string, options []types.SelectOption) *FormBuilder {
	field := types.Field{
		Name:    name,
		Type:    types.FieldTypeSelect,
		Label:   label,
		Options: options,
	}
	return fb.AddField(field)
}

// AddMultiSelectField добавляет поле множественного выбора
func (fb *FormBuilder) AddMultiSelectField(name, label string, options []types.SelectOption) *FormBuilder {
	field := types.Field{
		Name:     name,
		Type:     types.FieldTypeSelect,
		Label:    label,
		Options:  options,
		Multiple: true,
	}
	return fb.AddField(field)
}

// AddCheckboxField добавляет поле чекбокса
func (fb *FormBuilder) AddCheckboxField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypeCheckbox,
		Label: label,
	}
	return fb.AddField(field)
}

// AddTextareaField добавляет поле текстовой области
func (fb *FormBuilder) AddTextareaField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypeTextarea,
		Label: label,
	}
	return fb.AddField(field)
}

// AddDateField добавляет поле даты
func (fb *FormBuilder) AddDateField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypeDate,
		Label: label,
	}
	return fb.AddField(field)
}

// AddFileField добавляет поле файла
func (fb *FormBuilder) AddFileField(name, label string) *FormBuilder {
	field := types.Field{
		Name:  name,
		Type:  types.FieldTypeFile,
		Label: label,
	}
	return fb.AddField(field)
}

// AddHiddenField добавляет скрытое поле
func (fb *FormBuilder) AddHiddenField(name string, value interface{}) *FormBuilder {
	field := types.Field{
		Name:         name,
		Type:         types.FieldTypeHidden,
		DefaultValue: value,
	}
	return fb.AddField(field)
}

// AddTableField добавляет поле таблицы
func (fb *FormBuilder) AddTableField(name, label string) *TableFieldBuilder {
	field := types.Field{
		Name:        name,
		Type:        types.FieldTypeTable,
		Label:       label,
		TableConfig: &types.TableConfig{
			Columns:    make([]types.TableColumn, 0),
			Pagination: true,
			PageSize:   10,
			Sortable:   true,
			Filterable: true,
		},
	}
	
	return &TableFieldBuilder{
		field: &field,
	}
}

// AddGroup добавляет группу полей
func (fb *FormBuilder) AddGroup(name, title string, fields []string) *FormBuilder {
	group := types.FieldGroup{
		Name:   name,
		Title:  title,
		Fields: fields,
	}
	fb.form.Groups = append(fb.form.Groups, group)
	return fb
}

// OnPost устанавливает обработчик POST запросов
func (fb *FormBuilder) OnPost(handler types.FormHandler) *FormBuilder {
	fb.form.OnPost = handler
	return fb
}

// OnGet устанавливает обработчик GET запросов
func (fb *FormBuilder) OnGet(handler types.GetHandler) *FormBuilder {
	fb.form.OnGet = handler
	return fb
}

// Build завершает построение формы
func (fb *FormBuilder) Build() *types.Form {
	return fb.form
}

// TableFieldBuilder методы для настройки таблицы

// AddColumn добавляет колонку в таблицу
func (tfb *TableFieldBuilder) AddColumn(column types.TableColumn) *TableFieldBuilder {
	tfb.field.TableConfig.Columns = append(tfb.field.TableConfig.Columns, column)
	return tfb
}

// AddTextColumn добавляет текстовую колонку
func (tfb *TableFieldBuilder) AddTextColumn(key, title string) *TableFieldBuilder {
	column := types.TableColumn{
		Key:   key,
		Title: title,
		Type:  types.FieldTypeText,
	}
	return tfb.AddColumn(column)
}

// AddEmailColumn добавляет email колонку
func (tfb *TableFieldBuilder) AddEmailColumn(key, title string) *TableFieldBuilder {
	column := types.TableColumn{
		Key:   key,
		Title: title,
		Type:  types.FieldTypeEmail,
	}
	return tfb.AddColumn(column)
}

// AddNumberColumn добавляет числовую колонку
func (tfb *TableFieldBuilder) AddNumberColumn(key, title string) *TableFieldBuilder {
	column := types.TableColumn{
		Key:   key,
		Title: title,
		Type:  types.FieldTypeNumber,
	}
	return tfb.AddColumn(column)
}

// AddSelectColumn добавляет select колонку
func (tfb *TableFieldBuilder) AddSelectColumn(key, title string, options []types.SelectOption) *TableFieldBuilder {
	column := types.TableColumn{
		Key:     key,
		Title:   title,
		Type:    types.FieldTypeSelect,
		Options: options,
	}
	return tfb.AddColumn(column)
}

// AddMultiSelectColumn добавляет multiselect колонку
func (tfb *TableFieldBuilder) AddMultiSelectColumn(key, title string, options []types.SelectOption) *TableFieldBuilder {
	column := types.TableColumn{
		Key:      key,
		Title:    title,
		Type:     types.FieldTypeSelect,
		Options:  options,
		Multiple: true,
	}
	return tfb.AddColumn(column)
}

// AddCheckboxColumn добавляет checkbox колонку
func (tfb *TableFieldBuilder) AddCheckboxColumn(key, title string) *TableFieldBuilder {
	column := types.TableColumn{
		Key:   key,
		Title: title,
		Type:  types.FieldTypeCheckbox,
	}
	return tfb.AddColumn(column)
}

// AddDateColumn добавляет колонку даты
func (tfb *TableFieldBuilder) AddDateColumn(key, title string) *TableFieldBuilder {
	column := types.TableColumn{
		Key:   key,
		Title: title,
		Type:  types.FieldTypeDate,
	}
	return tfb.AddColumn(column)
}

// WithSortable делает колонку сортируемой
func (tfb *TableFieldBuilder) WithSortable() *TableFieldBuilder {
	if len(tfb.field.TableConfig.Columns) > 0 {
		lastIdx := len(tfb.field.TableConfig.Columns) - 1
		tfb.field.TableConfig.Columns[lastIdx].Sortable = true
	}
	return tfb
}

// WithFilterable делает колонку фильтруемой
func (tfb *TableFieldBuilder) WithFilterable() *TableFieldBuilder {
	if len(tfb.field.TableConfig.Columns) > 0 {
		lastIdx := len(tfb.field.TableConfig.Columns) - 1
		tfb.field.TableConfig.Columns[lastIdx].Filterable = true
	}
	return tfb
}

// WithWidth устанавливает ширину колонки
func (tfb *TableFieldBuilder) WithWidth(width string) *TableFieldBuilder {
	if len(tfb.field.TableConfig.Columns) > 0 {
		lastIdx := len(tfb.field.TableConfig.Columns) - 1
		tfb.field.TableConfig.Columns[lastIdx].Width = width
	}
	return tfb
}

// WithAlign устанавливает выравнивание колонки
func (tfb *TableFieldBuilder) WithAlign(align string) *TableFieldBuilder {
	if len(tfb.field.TableConfig.Columns) > 0 {
		lastIdx := len(tfb.field.TableConfig.Columns) - 1
		tfb.field.TableConfig.Columns[lastIdx].Align = align
	}
	return tfb
}

// WithPagination включает/выключает пагинацию
func (tfb *TableFieldBuilder) WithPagination(enabled bool) *TableFieldBuilder {
	tfb.field.TableConfig.Pagination = enabled
	return tfb
}

// WithPageSize устанавливает размер страницы
func (tfb *TableFieldBuilder) WithPageSize(size int) *TableFieldBuilder {
	tfb.field.TableConfig.PageSize = size
	return tfb
}

// WithSelectable делает строки выбираемыми
func (tfb *TableFieldBuilder) WithSelectable(enabled bool) *TableFieldBuilder {
	tfb.field.TableConfig.Selectable = enabled
	return tfb
}

// WithEditable делает строки редактируемыми
func (tfb *TableFieldBuilder) WithEditable(enabled bool) *TableFieldBuilder {
	tfb.field.TableConfig.Editable = enabled
	return tfb
}

// OnGet устанавливает обработчик получения данных таблицы
func (tfb *TableFieldBuilder) OnGet(handler types.TableHandler) *TableFieldBuilder {
	tfb.field.TableConfig.OnGet = handler
	return tfb
}

// Build завершает построение поля таблицы и возвращает FormBuilder
func (tfb *TableFieldBuilder) Build(fb *FormBuilder) *FormBuilder {
	return fb.AddField(*tfb.field)
}
