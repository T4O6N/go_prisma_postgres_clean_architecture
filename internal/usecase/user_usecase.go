package usecase

import (
	"context"
	"sample-project/internal/dto"
	"sample-project/internal/entity"
	"sample-project/internal/repository"
)

type UserUseCase interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, id int, updateUserDto dto.UpdateUserDto) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUseCase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetUsers(ctx context.Context) ([]entity.User, error) {
	return u.repo.GetAllUsers(ctx)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *userUsecase) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	return u.repo.CreateUser(ctx, user)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id int, updateUserDto dto.UpdateUserDto) (*entity.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateUserDto.Name != "" {
		user.Name = updateUserDto.Name
	}

	if updateUserDto.Email != "" {
		user.Email = updateUserDto.Email
	}

	if updateUserDto.SubjectID != 0 {
		user.SubjectID = updateUserDto.SubjectID
	}

	return u.repo.UpdateUser(ctx, id, *user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int) error {
	return u.repo.DeleteUser(ctx, id)
}
