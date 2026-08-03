package repository

import (
	"context"

	"github.com/M15z/restaurantops-api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m15z/restaurantops-api/internal/models"
)

type MenuItemRepository struct {
	pool *pgxpool.Pool
}

func NewMenuItemRepository(pool *pgxpool.Pool) *MenuItemRepository {
	return &MenuItemRepository{pool: pool}
}

func (r *MenuItemRepository) Create(ctx context.Context, m models.MenuItem) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `INSERT INTO menu_items (category_id, name, price, is_available)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		m.Category.ID, m.Name, m.Price, m.IsAvailable).Scan(&id)

	return id, err
}

func (r *MenuItemRepository) GetAll(ctx context.Context) ([]models.MenuItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT mi.id,mi.name, mi.price, mi.is_available, mi.created_at, c.id, c.name
	FROM menu_items mi
	JOIN categories c ON c.id = mi.category_id
	ORDER BY c.name, mi.name`)

	if err != nil {
		return nil, err
	}

	var items []models.MenuItem

	for rows.Next() {
		var m models.MenuItem
		if err := rows.Scan(&m.ID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt,
			&m.Category_ID, &m.Category.Name); err != nil {
			return nil, err
		}

		items = append(items, m)
	}

	return items, rows.Err()
}

func (r *MenuItemRepository) GetByID(ctx context.Context, id int64) (*models.MenuItem, error) {
	var m models.MenuItem

	err := r.pool.QueryRow(ctx, `SELECT mi.id,mi.name, mi.price, mi.is_available,
				mi.created_at, c.id, c.name
	FROM menu_items mi
	JOIN categories c ON c.id = mi.category_id
	WHERE mi.id = $1`, id).Scan(&m.ID, &m.Name, &m.Price, &m.IsAvailable, &m.CreatedAt,
		&m.Category.ID, &m.Category.Name)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *MenuItemRepository) Update(ctx context.Context, m models.MenuItem) error {

	_, err := r.pool.Exec(ctx, `UPDATE menu_items
				SET category_id = $1, name = $2, price = $3, is_available = $4
			WHERE id = $5`, m.Category.ID, m.Name, m.Price, m.IsAvailable, m.ID)
	return err
}

func (r *MenuItemRepository) SetAvailable(ctx context.Context, id int64, available bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE menu_items 
	SET is_available = $1 WHERE id = $2`, available, id)
	return err
}
