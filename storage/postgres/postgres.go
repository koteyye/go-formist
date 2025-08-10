package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koteyye/go-formist/storage"
)

// PostgresStorage реализация Storage для PostgreSQL
type PostgresStorage struct {
	pool *pgxpool.Pool
	sb   sq.StatementBuilderType
}

// NewPostgresStorage создает новое подключение к PostgreSQL
func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить DSN: %w", err)
	}

	// Настраиваем пул соединений
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул соединений: %w", err)
	}

	// Проверяем соединение
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("не удалось проверить соединение: %w", err)
	}

	ps := &PostgresStorage{
		pool: pool,
		sb:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
	
	// Создаем таблицу если её нет
	if err := ps.createTable(ctx); err != nil {
		return nil, fmt.Errorf("не удалось создать таблицу: %w", err)
	}

	return ps, nil
}

// createTable создает таблицу для хранения роутов
func (ps *PostgresStorage) createTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS formist_routes (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		path VARCHAR(255) NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		icon VARCHAR(100),
		type VARCHAR(50) NOT NULL CHECK (type IN ('form', 'page')),
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_routes_type ON formist_routes(type);
	CREATE INDEX IF NOT EXISTS idx_routes_name ON formist_routes(name);
	`

	_, err := ps.pool.Exec(ctx, query)
	return err
}

// SaveRoute сохраняет или обновляет роут
func (ps *PostgresStorage) SaveRoute(ctx context.Context, route *storage.Route) error {
	// Генерируем ID если его нет
	if route.ID == "" {
		route.ID = fmt.Sprintf("%s_%s_%d", route.Type, route.Name, time.Now().Unix())
	}

	// Устанавливаем временные метки
	now := time.Now()
	if route.CreatedAt.IsZero() {
		route.CreatedAt = now
	}
	route.UpdatedAt = now

	// Используем squirrel для построения запроса
	query, args, err := ps.sb.
		Insert("formist_routes").
		Columns("id", "name", "path", "title", "description", "icon", "type", "created_at", "updated_at").
		Values(
			route.ID,
			route.Name,
			route.Path,
			route.Title,
			route.Description,
			route.Icon,
			route.Type,
			route.CreatedAt,
			route.UpdatedAt,
		).
		Suffix(`
			ON CONFLICT (id) DO UPDATE SET
				name = EXCLUDED.name,
				path = EXCLUDED.path,
				title = EXCLUDED.title,
				description = EXCLUDED.description,
				icon = EXCLUDED.icon,
				type = EXCLUDED.type,
				updated_at = EXCLUDED.updated_at
		`).
		ToSql()

	if err != nil {
		return fmt.Errorf("не удалось построить запрос: %w", err)
	}

	_, err = ps.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось сохранить роут: %w", err)
	}

	return nil
}

// GetRoutes возвращает все роуты
func (ps *PostgresStorage) GetRoutes(ctx context.Context) ([]*storage.Route, error) {
	query, args, err := ps.sb.
		Select("id", "name", "path", "title", "description", "icon", "type", "created_at", "updated_at").
		From("formist_routes").
		OrderBy("type ASC", "title ASC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("не удалось построить запрос: %w", err)
	}

	rows, err := ps.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer rows.Close()

	routes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*storage.Route, error) {
		route := &storage.Route{}
		var description, icon *string
		
		err := row.Scan(
			&route.ID,
			&route.Name,
			&route.Path,
			&route.Title,
			&description,
			&icon,
			&route.Type,
			&route.CreatedAt,
			&route.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description != nil {
			route.Description = *description
		}
		if icon != nil {
			route.Icon = *icon
		}

		return route, nil
	})

	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать результаты: %w", err)
	}

	return routes, nil
}

// DeleteRoute удаляет роут по ID
func (ps *PostgresStorage) DeleteRoute(ctx context.Context, id string) error {
	query, args, err := ps.sb.
		Delete("formist_routes").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("не удалось построить запрос: %w", err)
	}

	result, err := ps.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось удалить роут: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("роут с ID %s не найден", id)
	}

	return nil
}

// Close закрывает пул соединений
func (ps *PostgresStorage) Close() error {
	ps.pool.Close()
	return nil
}