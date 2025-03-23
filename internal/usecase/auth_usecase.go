package usecase

import (
	"context"
	"errors"
	"sample-project/internal/entity"
	"sample-project/internal/repository"
	"sample-project/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(ctx context.Context, name, password string) (string, string, error)
	GetUserProfile(ctx context.Context, token string) (*entity.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUseCase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) Login(ctx context.Context, name, password string) (string, string, error) {
	user, err := u.userRepo.GetUserByName(ctx, name)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := utils.GenerateToken(user.ID, user.Name)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *authUsecase) GetUserProfile(ctx context.Context, token string) (*entity.User, error) {
	claims, err := utils.ValidateToken(token, false)
	if err != nil {
		return nil, errors.New("invalid token or expired token")
	}

	user, err := u.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
