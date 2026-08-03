package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/m15z/restaurantops-api/internal/models"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u models.User) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		u.Name, u.Email, u.PasswordHash, u.Role,
	).Scan(&id)
	return id, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx, `SELECT id, name, email, password_hash, role, created_at
	FROM users
	WHERE email = $1`, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil

}
