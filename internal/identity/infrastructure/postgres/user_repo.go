package postgres

import (
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
