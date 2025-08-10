package storage

import (
	"context"
	"time"
)

// Route представляет роут для навигации
type Route struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Path        string    `json:"path" db:"path"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Icon        string    `json:"icon,omitempty" db:"icon"`
	Type        string    `json:"type" db:"type"` // form или page
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Storage интерфейс для хранения роутов
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