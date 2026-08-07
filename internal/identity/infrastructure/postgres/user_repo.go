package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

type userRow struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}

func (r userRow) toDomain() (*domain.User, error) {
	email, err := domain.NewEmail(r.Email)
	if err != nil {
		return nil, err
	}
	return domain.NewUserFromRepository(r.ID, email, r.PasswordHash, domain.Role(r.Role), r.CreatedAt), nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {

	const query = `INSERT INTO users (id, email, password_hash, role, created_at)
			VALUES (:id, :email, :password_hash, :role, :created_at)
			ON CONFLICT (id) DO UPDATE SET
				email = EXCLUDED.email,
				password_hash = EXCLUDED.password_hash,
				role = EXCLUDED.role`
	row := userRow{
		ID:           user.ID(),
		Email:        user.Email().String(),
		PasswordHash: user.PasswordHash(),
		Role:         string(user.Role()),
		CreatedAt:    user.CreatedAt(),
	}
	_, err := r.db.NamedExecContext(ctx, query, row)
	return err

}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1`
	var row userRow
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return row.toDomain()
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	const query = `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	var row userRow
	if err := r.db.GetContext(ctx, &row, query, email.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return row.toDomain()
}
