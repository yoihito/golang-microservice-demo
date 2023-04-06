package repositories

import (
	"auth/internal/entities"
	"context"
)

type User struct {
	db Datastore
}

func NewUser(db Datastore) *User {
	return &User{db}
}

func (u *User) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}
	if err := u.db.QueryRowContext(
		ctx,
		"SELECT email, password_digest FROM users WHERE email = $1",
		email).Scan(&user.Email, &user.PasswordDigest); err != nil {
		return nil, err
	}

	return user, nil
}
