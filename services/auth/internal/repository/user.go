package repository

import (
	"context"
	"copo/auth/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *model.User, hashedPassword string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, email, password, name, role)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
	`, u.Email, hashedPassword, u.Name, u.Role)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password, name, role, created_at
		FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) UpdateName(ctx context.Context, email, name string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx, `
		UPDATE users SET name = $1
		WHERE email = $2
		RETURNING id, email, name, role, created_at
	`, name, email).Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt)
	return u, err
}
