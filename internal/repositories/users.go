package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type User struct {
	ID        int       `redis:"id"`
	Name      string    `redis:"name"`
	Email     string    `redis:"email"`
	Password  string    `redis:"password"`
	Role      string    `redis:"role"`
	Status    string    `redis:"status"`
	CreatedAt time.Time `redis:"created_at"`
	UpdatedAt time.Time `redis:"updated_at"`
}

// CreateUser creates a user with given user model in parameter
func (r *Repository) CreateUser(ctx context.Context, userRequest User) error {
	user := models.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
		Role:     userRequest.Role,
		Status:   userRequest.Status,
	}
	if err := user.Insert(ctx, boil.GetContextDB(), boil.Infer()); err != nil {
		return err
	}
	return nil
}

// GetUser gets a user from db by user id
func (r *Repository) GetUser(ctx context.Context, id int) (models.User, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("user:%d", id))
	if len(res.Val()) == 0 {
		user, errDb := models.FindUser(ctx, boil.GetContextDB(), id)
		if errDb != nil {
			if errors.Is(errDb, sql.ErrNoRows) {
				return models.User{}, ErrUserNotFound
			}
			return models.User{}, errDb
		}

		userCache := User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		// set cache
		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("user:%d", id), userCache); errCache.Err() != nil {
			return models.User{}, errCache.Err()
		}

		return *user, nil
	}

	var userScan User
	if err := res.Scan(&userScan); err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:        userScan.ID,
		Name:      userScan.Name,
		Email:     userScan.Email,
		Password:  userScan.Password,
		Role:      userScan.Role,
		Status:    userScan.Status,
		CreatedAt: userScan.CreatedAt,
		UpdatedAt: userScan.UpdatedAt,
	}, nil
}

const (
	userIDDefault       = 1
	userNameDefault     = "user default"
	userEmailDefault    = "user@gmail.com"
	userPasswordDefault = "userpassword"
)

// GetUserDefault gets the id of default user
func (r *Repository) GetUserDefault(ctx context.Context) (models.User, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("user:%d", userIDDefault))
	if len(res.Val()) == 0 {
		user, err := models.FindUser(ctx, boil.GetContextDB(), userIDDefault)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				userDefault := models.User{
					ID:       userIDDefault,
					Name:     userNameDefault,
					Email:    userEmailDefault,
					Password: userPasswordDefault,
					Role:     "user",
					Status:   "activated",
				}

				if errQuery := userDefault.Insert(ctx, boil.GetContextDB(), boil.Infer()); errQuery != nil {
					return models.User{}, errQuery
				}

				userCache := User{
					ID:        userDefault.ID,
					Name:      userDefault.Name,
					Email:     userDefault.Email,
					Password:  userDefault.Password,
					Role:      userDefault.Role,
					Status:    userDefault.Status,
					CreatedAt: userDefault.CreatedAt,
					UpdatedAt: userDefault.UpdatedAt,
				}

				// set cache
				if errCache := r.Redis.HSet(ctx, fmt.Sprintf("user:%d", userDefault.ID), userCache); errCache.Err() != nil {
					return models.User{}, errCache.Err()
				}

				return userDefault, nil
			}
			return models.User{}, err
		}

		userCache := User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		// set cache
		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("user:%d", user.ID), userCache); errCache.Err() != nil {
			return models.User{}, errCache.Err()
		}
		return *user, nil
	}

	var userScan User
	if err := res.Scan(&userScan); err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:       userScan.ID,
		Name:     userScan.Name,
		Email:    userScan.Email,
		Password: userScan.Password,
		Role:     userScan.Role,
		Status:   userScan.Status,
	}, nil
}
