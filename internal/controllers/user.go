package controllers

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/botbooker/botbooker/internal/database"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	dbConn *bun.DB
}

func UserControllerInitializer(conn *bun.DB) *UserController {
	return &UserController{
		dbConn: conn,
	}
}

func hashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password is empty")
	}
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters long")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(hashedPassword), nil
}

func (u *UserController) CreateUser(user *database.User) (bool, error) {
	ctx := context.Background()
	if user.PasswordHash == "" {
		user.PasswordHash = os.Getenv("DEFAULT_PASSWORD")
	}
	hashedPassword, err := hashPassword(user.PasswordHash)
	if err != nil {
		return false, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)
	_, err = u.dbConn.NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to create user: %w", err)
	}
	return true, nil
}
